metadata :name        => "emulator",
         :description => "choria-emulator manager agent",
         :author      => "R.I.Pienaar <rip@devco.net>",
         :license     => "Apache-2.0",
         :version     => "0.0.1",
         :url         => "http://choria.io",
         :timeout     => 60

requires :mcollective => "2.9.0"

action "status", :description => "Status of the running emulator" do
  input :port,
        :description => "Port to query for status",
        :prompt => "Monitor Port",
        :type => :integer,
        :optional => true,
        :default => 8080

  output :name,
         :description => "Instance name",
         :display_as => "Name"

  output :running,
         :description => "Is the instance running",
         :display_as => "Running"

  output :pid,
         :description => "Running PID",
         :display_as => "PID"

  output :tls,
         :description => "TLS Enabled",
         :display_as => "TLS"

  output :memory,
         :description => "Memory used in bytes",
         :display_as => "Memory (B)"

  summarize do
    aggregate summary(:running)
    aggregate summary(:tls)
  end
end

action "download", :description => "Downloads the emulator binary" do
  input :http,
        :description => "HTTP or HTTPS URL to fetch the file from",
        :prompt => "Emulator source URL",
        :type => :string,
        :optional => false

  input :target,
        :description => "Location to store the downloaded emulator in",
        :prompt => "Target",
        :type => :string,
        :optional => false

  output :status,
         :description => "true if the emulator was downloaded",
         :display_as => "Downloaded"

  summarize do
    aggregate summary(:status)
  end
end

action "stop", :description => "Stops the running emulator instance" do
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
        :optional => true

  input :agents,
        :prompt => "Agents",
        :description => "Number of emulated* agents the emulator will host",
        :type => :integer,
        :optional => true,
        :default => 1

  input :collectives,
        :prompt => "Subcollectives",
        :description => "Number of subcollective the emulator will join",
        :type => :integer,
        :optional => true,
        :default => 1

  input :servers,
        :prompt => "Servers to connect to",
        :description => "Comma separated list of host:port pairs",
        :type => :string,
        :optional => true

  input :emulator,
        :prompt => "Path to emulator binary",
        :description => "Where to find the emulator binary",
        :type => :string,
        :optional => false

  input :monitor,
        :description => "Port to listen for monitoring requests",
        :prompt => "Monitor Port",
        :type => :integer,
        :optional => true,
        :default => 8080

  output :status,
         :description => "true if the emulator started",
         :display_as => "Started"

  summarize do
    aggregate summary(:status)
  end
end
