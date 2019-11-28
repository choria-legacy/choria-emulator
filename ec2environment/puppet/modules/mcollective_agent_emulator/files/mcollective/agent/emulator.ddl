metadata :name        => "emulator",
         :description => "choria-emulator manager agent",
         :author      => "R.I.Pienaar <rip@devco.net>",
         :license     => "Apache-2.0",
         :version     => "0.0.1",
         :url         => "http://choria.io",
         :timeout     => 60


action "download", :description => "Downloads the emulator binary" do
  display :failed

  input :http,
        :prompt      => "Source URL",
        :description => "HTTP or HTTPS URL to fetch the file from",
        :type        => :string,
        :validation  => '.',
        :maxlength   => 256,
        :optional    => false


  input :target,
        :prompt      => "Downloaded file to be stored here",
        :description => "Name under /tmp/choria-emulator",
        :type        => :string,
        :validation  => '^[a-zA-Z\d\.-]+$',
        :maxlength   => 256,
        :optional    => false




  output :md5,
         :description => "md5 hash of downloaded file",
         :display_as  => "MD5"

  output :size,
         :description => "Size of file downloaded",
         :display_as  => "Size"

  output :success,
         :description => "true if the file was downloaded successfully",
         :display_as  => "Downloaded"

  summarize do
    aggregate summary(:success)
    aggregate summary(:size)
    aggregate summary(:md5)
  end
end

action "status", :description => "Status of the running emulator" do
  display :always

  input :port,
        :prompt      => "Monitor Port",
        :description => "Port to query for status",
        :type        => :integer,
        :default     => 8080,
        :optional    => false




  output :emulator,
         :description => "md5 hash of emulator binary",
         :display_as  => "Emulator"

  output :memory,
         :description => "Memory used in bytes",
         :display_as  => "Memory (B)"

  output :name,
         :description => "Instance name",
         :default     => nil,
         :display_as  => "Name"

  output :pid,
         :description => "Running PID",
         :default     => nil,
         :display_as  => "PID"

  output :running,
         :description => "Is the instance running",
         :display_as  => "Running"

  output :tls,
         :description => "TLS Enabled",
         :display_as  => "TLS"

  summarize do
    aggregate summary(:running)
    aggregate summary(:tls)
    aggregate summary(:emulator)
  end
end

action "start", :description => "Start an emulator instance" do
  display :failed

  input :agents,
        :prompt      => "Agents",
        :description => "Number of emulated* agents the emulator will host",
        :type        => :integer,
        :default     => 1,
        :optional    => false


  input :collectives,
        :prompt      => "Subcollectives",
        :description => "Number of subcollective the emulator will join",
        :type        => :integer,
        :default     => 1,
        :optional    => false


  input :credentials,
        :prompt      => "Credentials",
        :description => "Base64 encoded credentials to use when connecting to NATS",
        :type        => :string,
        :validation  => '.',
        :maxlength   => 2048,
        :optional    => true


  input :instances,
        :prompt      => "Instances",
        :description => "Number of simulated choria instances the emulator will host",
        :type        => :integer,
        :default     => 1,
        :optional    => false


  input :monitor,
        :prompt      => "Monitor Port",
        :description => "Port to listen for monitoring requests",
        :type        => :integer,
        :default     => 8080,
        :optional    => false


  input :name,
        :prompt      => "Name",
        :description => "Instance Name",
        :type        => :string,
        :validation  => '^[\w\-_]+$',
        :maxlength   => 16,
        :optional    => true


  input :servers,
        :prompt      => "Servers to connect to",
        :description => "Comma separated list of host:port pairs",
        :type        => :string,
        :validation  => '.',
        :maxlength   => 256,
        :optional    => true


  input :tls,
        :prompt      => "TLS",
        :description => "Run with TLS enabled",
        :type        => :boolean,
        :optional    => false




  output :status,
         :description => "true if the emulator started",
         :display_as  => "Started"

  summarize do
    aggregate summary(:status)
  end
end

action "start_nats", :description => "Start a local NATS instance" do
  display :failed

  input :monitor_port,
        :prompt      => "Monitor Port",
        :description => "The port to listen on for monitoring requests",
        :type        => :integer,
        :default     => 8222,
        :optional    => false


  input :port,
        :prompt      => "Client Port",
        :description => "Port for clients to connect to",
        :type        => :integer,
        :default     => 4222,
        :optional    => false




  output :running,
         :description => "true if the NATS started",
         :display_as  => "Started"

  summarize do
    aggregate summary(:running)
  end
end

action "stop", :description => "Stops the running emulator instance" do
  display :failed

  input :port,
        :prompt      => "Monitor Port",
        :description => "Port to query for status",
        :type        => :integer,
        :default     => 8080,
        :optional    => false




  output :status,
         :description => "true if the emulator stopped",
         :display_as  => "Stopped"

  summarize do
    aggregate summary(:status)
  end
end

action "stop_nats", :description => "Stops a running NATS instance" do
  display :failed



  output :stopped,
         :description => "true if the NATS stopped",
         :display_as  => "Stopped"

  summarize do
    aggregate summary(:stopped)
  end
end

