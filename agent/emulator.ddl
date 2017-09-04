metadata :name        => "emulator",
         :description => "choria-emulator manager agent",
         :author      => "R.I.Pienaar <rip@devco.net>",
         :license     => "Apache-2.0",
         :version     => "0.0.1",
         :url         => "http://choria.io",
         :timeout     => 60

requires :mcollective => "2.9.0"

action "status", :description => "Status of the running emulator" do
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
        :optional => false,
        :maxlength => "256",
        :validation => "."

  output :success,
         :description => "true if the emulator was downloaded",
         :display_as => "Downloaded"

  output :size,
         :description => "Size of file downloaded",
         :display_as => "Size"

  summarize do
    aggregate summary(:success)
    aggregate summary(:size)
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
        :validation => '^\w+$'

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
        :maxlength => "256",
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

  output :status,
         :description => "true if the emulator started",
         :display_as => "Started"

  summarize do
    aggregate summary(:status)
  end
end
