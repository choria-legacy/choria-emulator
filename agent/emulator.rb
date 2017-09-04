require 'net/http'

module MCollective
  module Agent
    class Emulator<RPC::Agent
      action "download" do
        reply[:success] = false

        FileUtils.mkdir_p("/tmp/choria-emulator")

        if request[:http]
          begin
            download_http(request[:http], "/tmp/choria-emulator/choria-emulator")
          rescue
            reply.fail!("Failed to download %s: %s: %s" % [request[:http], $!.class, $!.to_s])
          end

          FileUtils.chmod(0755, "/tmp/choria-emulator/choria-emulator")

          stat = File::Stat.new("/tmp/choria-emulator/choria-emulator")
          reply[:size] = stat.size
          reply[:success] = true
        else
          reply.fail("No valid download location given")
        end
      end

      action "start" do
        unless File.executable?("/tmp/choria-emulator/choria-emulator")
          reply.fail!("Cannot start /tmp/choria-emulator/choria-emulator does not exist or is not executable")
        end

        if up?(request[:monitor])
          reply.fail!("Cannot start, emulator is already running")
        end

        args = []

        args << "--name %s" % request[:name]
        args << "--instances %d" % request[:instances]
        args << "--http-port %d" % request[:monitor]
        args << "--config /dev/null"

        args << "--agents %d" % request[:agents] if request[:agents]
        args << "--collectives %d" % request[:collectives] if request[:collectives]
        args << "--server %s" % request[:servers].gsub(" ", "") if request[:servers]
        args << "--tls" if request[:tls]

        out = []
        err = []
        Log.info("Running: %s" % args.join(" "))

        run('(/tmp/choria-emulator/choria-emulator %s 2>&1 >> /tmp/choria-emulator/log &) &' % args.join(" "), :stdout => out, :stderr => err)

        sleep 1

        reply[:status] = up?(request[:monitor])
      end

      action "status" do
        if down?(request[:port])
          reply[:running] = false
          break
        end

        status = emulator_status(request[:port])

        reply[:name] = status["name"]
        reply[:running] = true
        reply[:pid] = status["config"]["pid"]
        reply[:tls] = status["config"]["TLS"] == 1
        reply[:memory] = status["memstats"]["Sys"]
      end

      action "stop" do
        reply[:status] = false

        if up?(request[:port])
          pid = emulator_pid(request[:port])

          reply.fail!("Could not determine PID for running emulator") unless pid

          Process.kill("HUP", pid)
          sleep 1
          reply[:status] = down?(request[:port])
        end
      end

      def up?(port)
        Log.debug(emulator_status(port).inspect)
        emulator_status(port)["status"] == "up"
      rescue
        Log.warn("%s: %s" % [$!.class, $!.to_s])
        false
      end

      def down?(port)
        !up?(port)
      end

      def emulator_pid(port)
        emulator_status(port)["config"]["pid"]
      end

      def emulator_status(port=8080)
        uri = URI.parse("http://localhost:%d/debug/vars" % port)
        response = Net::HTTP.get_response(uri)

        out = {
          "status" => "up",
          "code" => response.code
        }

        if response.code == "200"
          Log.debug(response.body)
          out.merge!(JSON.parse(response.body))
        end

        out
      end

      def download_http(url, target)
        uri = URI.parse(url)
        out = File.open(target, "wb")

        begin
          Net::HTTP.get_response(uri) do |resp|
            resp.read_body do |segment|
              out.write(segment)
            end
          end
        ensure
          out.close
        end
      end
    end
  end
end
