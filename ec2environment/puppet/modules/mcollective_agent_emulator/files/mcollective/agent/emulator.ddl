metadata :name        => "emulator",
         :description => "choria-emulator manager agent",
         :author      => "R.I.Pienaar <rip@devco.net>",
         :license     => "Apache-2.0",
         :version     => "0.0.1",
         :url         => "http://choria.io",
         :timeout     => 60

requires :mcollective => "2.9.0"

action "emulator_status", :description => "Status of the running emulator" do
  display :always

  input :port,
        :description => "Port to query for status",
        :prompt => "Monitor Port",
        :type => :integer,
        :optional => false,
        :default => 8080

  output :name,
         :description => "Instance name",
         :display_as => "Name",
         :default => "unknown"

  output :running,
         :description => "Is the instance running",
         :display_as => "Running",
         :default => false

  output :pid,
         :description => "Running PID",
         :display_as => "PID",
         :default => -1

  output :tls,
         :description => "TLS Enabled",
         :display_as => "TLS",
         :default => false

  output :memory,
         :description => "Memory used in bytes",
         :display_as => "Memory (B)",
         :default => 0

  output :emulator,
         :description => "md5 hash of emulator binary",
         :display_as => "Emulator",
         :default => ""

  summarize do
    aggregate summary(:running)
    aggregate summary(:tls)
    aggregate summary(:emulator)
  end
end

action "download", :description => "Downloads the emulator binary" do
  input :http,
        :description => "HTTP or HTTPS URL to fetch the file from",
        :prompt => "Source URL",
        :type => :string,
        :optional => false,
        :maxlength => "256",
        :validation => "."

  input :target,
        :description => "Name under /tmp/choria-emulator",
        :prompt => "Downloaded file to be stored here",
        :type => :string,
        :optional => false,
        :maxlength => "256",
        :validation => '^[a-zA-Z\d\.-]+$'

  output :success,
         :description => "true if the file was downloaded successfully",
         :display_as => "Downloaded"

  output :size,
         :description => "Size of file downloaded",
         :display_as => "Size"

  output :md5,
         :description => "md5 hash of downloaded file",
         :display_as => "MD5",
         :default => ""

  summarize do
    aggregate summary(:success)
    aggregate summary(:size)
    aggregate summary(:md5)
  end
end

action "stop", :description => "Stops the running emulator instance" do
  input :port,
        :description => "Port to query for status",
        :prompt => "Monitor Port",
        :type => :integer,
        :optional => false,
        :default => 8080

  output :status,
         :description => "true if the emulator stopped",
         :display_as => "Stopped"

  summarize do
    aggregate summary(:status)
  end
end

action "start", :description => "Start an emulator instance" do
  input :name,
        :prompt => "Name",
        :description => "Instance Name",
        :type => :string,
        :optional => true,
        :maxlength => "16",
        :validation => '^[\w\-_]+$'

  input :instances,
        :prompt => "Instances",
        :description => "Number of simulated choria instances the emulator will host",
        :type => :integer,
        :optional => false,
        :default => 1

  input :agents,
        :prompt => "Agents",
        :description => "Number of emulated* agents the emulator will host",
        :type => :integer,
        :optional => false,
        :default => 1

  input :collectives,
        :prompt => "Subcollectives",
        :description => "Number of subcollective the emulator will join",
        :type => :integer,
        :optional => false,
        :default => 1

  input :servers,
        :prompt => "Servers to connect to",
        :description => "Comma separated list of host:port pairs",
        :type => :string,
        :maxlength => 256,
        :optional => true,
        :validation => "."

  input :monitor,
        :description => "Port to listen for monitoring requests",
        :prompt => "Monitor Port",
        :type => :integer,
        :optional => false,
        :default => 8080

  input :tls,
        :description => "Run with TLS enabled",
        :prompt => "TLS",
        :type => :boolean,
        :optional => false,
        :default => false

  input :credentials,
        :description => "Base64 encoded credentials to use when connecting to NATS",
        :prompt => "Credentials",
        :type => :string,
        :optional => true,
        :maxlength => 2048,
        :validation => "."

  output :status,
         :description => "true if the emulator started",
         :display_as => "Started"

  summarize do
    aggregate summary(:status)
  end
end

action "start_nats", :description => "Start a local NATS instance" do
  input :port,
        :prompt => "Client Port",
        :description => "Port for clients to connect to",
        :default => 4222,
        :type => :integer,
        :optional => false

  input :monitor_port,
        :prompt => "Monitor Port",
        :description => 8222,
        :type => :integer,
        :optional => false

  output :running,
         :description => "true if the NATS started",
         :display_as => "Started"

  summarize do
    aggregate summary(:running)
  end
end

action "stop_nats", :description => "Stops a running NATS instance" do
  output :stopped,
         :description => "true if the NATS stopped",
         :display_as => "Stopped"

  summarize do
    aggregate summary(:stopped)
  end
end

action "start_federation", :description => "Start a local Federation Broker instance" do
  input :name,
        :description => "The Collective Name this broker will serve",
        :prompt => "Collective Name",
        :type => :string,
        :maxlength => 50,
        :optional => true,
        :validation => '.+'

  input :tls,
        :description => "Run with TLS enabled",
        :prompt => "TLS",
        :type => :boolean,
        :optional => false,
        :default => false

  input :federation_servers,
        :description => "Federation NATS Servers, comma separated list of host:port pairs",
        :prompt => "Federation Servers",
        :type => :string,
        :optional => false,
        :maxlength => 256,
        :validation => '.+'

  input :collective_servers,
        :description => "Collective NATS Servers, comma separated list of host:port pairs",
        :prompt => "Collective Servers",
        :type => :string,
        :optional => false,
        :maxlength => 256,
        :validation => '.+',
        :default => "localhost:4222"

  output :running,
         :description => "true if the Federation Broker started",
         :display_as => "Started"

  summarize do
    aggregate summary(:running)
  end
end

action "stop_federation", :description => "Stops a running Federation Broker instance" do
  output :stopped,
         :description => "true if the Federation Broker stopped",
         :display_as => "Stopped"

  summarize do
    aggregate summary(:stopped)
  end
end

