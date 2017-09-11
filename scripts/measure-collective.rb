#!/opt/puppetlabs/puppet/bin/ruby

require 'mcollective'
require 'csv'

include MCollective::RPC

@config = {
  :count => 10,
  :size => 10,
  :out => "reports",
  :stats_dir => nil
}

@stats = {}

module Enumerable
  def sum
    self.inject(0){|accum, i| accum + i }
  end

  def mean
    self.sum / self.length.to_f
  end

  def sample_variance
    m = self.mean
    sum = self.inject(0){|accum, i| accum +(i-m)**2 }
    sum/(self.length - 1).to_f
  end

  def standard_deviation
    return Math.sqrt(self.sample_variance)
  end

  def skewness
    return if length == 0
    return 0 if length == 1

    sum_cubed_deviation / ((length - 1) * cubed_standard_deviation.to_f)
  end

  def sum_cubed_deviation
    precalculated_mean = mean

    inject(0) {|sum, value| sum + (value - precalculated_mean) ** 3}
  end

  def cubed_standard_deviation
    standard_deviation ** 3
  end
end

def parse_cli
  oparser = MCollective::Optionparser.new({:verbose => false, :progress_bar => false}, nil, "common")
  options = oparser.parse do |parser, opts|
    parser.on("--disable-tls", "Disables TLS on the NATS connection") do
      $choria_unsafe_disable_nats_tls = true
    end

    parser.on("--discovery-timeout SECONDS", "--dt", Integer, "Discovery Timeout") do |v|
      options[:disctimeout] = v
    end

    parser.on("--count TESTS", "-c", Integer, "Number of tests to perform") do |v|
      @config[:count] = v
    end

    parser.on("--size SIZE", "-s", Integer, "Message payload size") do |v|
      @config[:size] = v
    end

    parser.on("--out DIR", "-o", "Directory to store results in") do |v|
      @config[:out] = v
    end

    parser.on("--description DESCRIPTION", String, "Test scenario description") do |v|
      @config[:description] = v
    end
  end

  abort("Please specify an output dir with --out") unless @config[:out]
  abort("Please specify a test description with --description") unless @config[:description]

  @config[:stats_dir] = File.join(@config[:out], Time.now.strftime("%Y%m%d-%H%M%S"))

  options
end

def setup_stats
  test_dir = @config[:stats_dir]

  abort("Output directory %s already exist" % test_dir) if File.exist?(test_dir)

  FileUtils.mkdir_p(test_dir)

  @stats[:outcomes] = CSV.open(File.join(test_dir, "outcomes.csv"), "w")
  @stats[:outcomes] << ["Success", "Responses", "NoResponses", "OKCount", "FailedCount", "Blocktime", "Overhead", "Totaltime", "ResponseMin", "ResponseMax", "ResponseAverage", "ResponseDeviation", "ResponseSkew"]
  @stats[:node_times] = CSV.open(File.join(test_dir, "node_times.csv"), "w")
end

def close_stats
  @stats[:outcomes].close
  @stats[:node_times].close
end

def run_test(options)
  client = rpcclient("emulated0", :options => options)
  found = client.discover.size

  abort("Did not discover any nodes") unless found > 0

  puts "Performing %d tests against %d nodes with a payload of %d bytes stats in %s" % [@config[:count], found, @config[:size], @config[:stats_dir]]

  begin
    @config[:count].times do |ctr|
      start_time = Time.now

      node_times = []

      client.generate(:size => @config[:size]) do |r|
        node_times << Time.now - start_time
      end

      end_time = Time.now

      @stats[:node_times] << node_times
      @stats[:outcomes] << [
        (client.stats.okcount == found ? 1 : 0),
        client.stats.responsesfrom.size,
        client.stats.noresponsefrom.size,
        client.stats.okcount,
        client.stats.failcount,
        client.stats.blocktime,
        end_time.to_f - start_time.to_f - client.stats.blocktime,
        end_time.to_f - start_time.to_f,
        node_times.min,
        node_times.max,
        node_times.mean,
        node_times.standard_deviation,
        node_times.skewness
      ]

      ok = client.stats.okcount == found ? MCollective::Util.colorize(:green, "✓") : MCollective::Util.colorize(:red, "X")
      puts " %d: REQUEST #: %-4d OK: %d / %d BLOCK TIME: %.3f ELAPSED TIME: %.3f %s" % [Time.now.to_i, ctr + 1, client.stats.okcount, found, client.stats.blocktime, end_time - client.stats.starttime, ok]
    end
  rescue
    p $!
  ensure
    close_stats
  end
end

def calculate_time_buckets(max)
  nbuckets = (max / 0.1).round + 1
  buckets = CSV.open(File.join(@config[:stats_dir], "time_bucketed.csv"), "w")
  buckets << nbuckets.times.map {|i| (i * 0.1).round(1)}

  CSV.foreach(File.join(@config[:stats_dir], "node_times.csv")) do |r|
    dat = {}
    nbuckets.times {|i| dat[(i * 0.1).round(1)] = 0}

    r.each do |response|
      dat[Float(response).round(1)] += 1
    end

    buckets << dat.values
  end

  buckets.close
end

def report
  ok = 0
  failed = 0
  min = +1.0/0.0
  max = -1.0/0.0
  avg = 0
  times = []

  CSV.foreach(File.join(@config[:stats_dir], "outcomes.csv"), :headers => true) do |r|
    r["Success"] == "1" ? ok += 1 : failed += 1

    t = Float(r["Totaltime"])
    min = t if min > t
    max = t if max < t
    times << t
  end

  puts
  puts "OK: %d FAILED: %d MIN: %.4f MAX: %.4f AVG: %.4f STDDEV: %.4f" % [ok, failed, min, max, times.mean, times.standard_deviation]

  calculate_time_buckets(max)

  File.open(File.join(@config[:stats_dir], "desciption.json"), "w") do |f|
    f.puts(@config.to_json)
  end
end

options = parse_cli

setup_stats

run_test(options)

report
