#!/bin/env ruby

require "net/http"
require "json"
require "pp"
require "optparse"
require "csv"

def url
  "http://%s:%s/%s" % [@config[:host], @config[:port], @config[:source]]
end

def get
  response = Net::HTTP.get_response(URI.parse(url))

  if response.code == "200"
    return JSON.parse(response.body)
  else
    raise("Received code %s from %s: %s" % [response.code, url, response.body])
  end
end

@config = {
  :host => "localhost",
  :port => 8222,
  :source => "subsz",
  :variable => "num_subscriptions",
  :out => nil
}

opt = OptionParser.new

opt.on("--host HOST", "Host where NATS runs (%s)" % @config[:host]) do |v|
  @config[:host] = v
end

opt.on("--port PORT", Integer, "NATS monitoring port (%s)" % @config[:port]) do |v|
  @config[:port] = v
end

opt.on("--source SOURCE", "NATS monitor url (%s)" % @config[:source]) do |v|
  @config[:source] = v
end

opt.on("--variable VAR", "NATS monitor variable (%s)" % @config[:variable]) do |v|
  @config[:variable] = v
end

opt.on("--out OUT", "File to write results in" % @config[:variable]) do |v|
  @config[:out] = v
end

opt.parse!

abort("Please specify a target file with --out") unless @config[:out]

if (start = get[@config[:variable]]) > 0
  puts "Waiting for %s to go below current threshold of %d" % [@config[:variable], start]

  loop do
    begin
      break if get[@config[:variable]] < start
      print "."
    rescue
      print "!"
    end
  end

  puts
end

start_time = nil
start_count = get[@config[:variable]]
previous = start_count
last_change = Time.at(0)
buckets = Hash.new(0)

puts "Waiting for /%s %s, starting at %d" % [@config[:source], @config[:variable], start_count]

loop do
  current = get[@config[:variable]]

  if current > start_count && start_time.nil?
    start_time = Time.now
    print ">"
  end

  if start_time
    if current > previous
      last_change = Time.now

      bucket = (last_change - start_time).round(1)

      buckets[bucket] ||= 0
      buckets[bucket] += current - previous

      previous = current

      print "."
    elsif Time.now - last_change > 10
      puts "<"
      puts "Stable on %d %s @ %s elapsed %.2fs" % [current, @config[:variable], last_change, last_change - start_time]
      break
    end
  end
end

out = CSV.open(@config[:out], "w")
out << buckets.keys
out << buckets.values
out.close
