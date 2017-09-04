metadata    :name        => "emulated0",
            :description => "Choria Agent emulated by choria-emulator",
            :author      => "R.I.Pienaar <rip@devco.net>",
            :license     => "Apache-2.0",
            :version     => "0.0.1",
            :url         => "http://choria.io",
            :timeout     => 120

requires :mcollective => "2.9.0"

action "generate", :description => "Generates random data of a given size" do
    input :size,
          :prompt => "Size",
          :description => "Amount of text to generate",
          :type => :integer,
          :optional => true,
          :default => 20

    output :message,
           :description => "Generated Message",
           :display_as => "Message"
end
