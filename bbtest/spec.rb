require 'turnip/rspec'
require 'json'
require 'thread'

Thread.abort_on_exception = true

RSpec.configure do |config|
  config.raise_error_for_unimplemented_steps = true
  config.color = true

  Dir.glob("./helpers/*_helper.rb") { |f| load f }
  config.include EventuallyHelper, :type => :feature
  Dir.glob("./steps/*_steps.rb") { |f| load f, true }

  config.before(:suite) do |_|
    print "[ suite starting ]\n"

    $call_id = 0

    # fixme input validation test that binary was built
    %x(ln -s /opt/binaries/linux-latest /bin/contract)

    ["/reports"].each { |folder|
      FileUtils.mkdir_p folder
      FileUtils.rm_rf Dir.glob("#{folder}/*")
    }

    print "[ suite started  ]\n"
  end

  config.after(:suite) do |_|
    print "\n[ suite ending   ]\n"

    get_containers = lambda do |image|
      containers = %x(docker ps -a | awk '{ print $1,$2 }' | grep #{image} | awk '{print $1 }' 2>/dev/null)
      return ($? == 0 ? containers.split("\n") : [])
    end

    teardown_container = lambda do |container|
      label = %x(docker inspect --format='{{.Name}}' #{container})
      label = ($? == 0 ? label.strip : container)

      %x(docker kill --signal="TERM" #{container} >/dev/null 2>&1 || :)
      %x(docker logs #{container} >/reports/#{label}.log 2>&1)
      %x(docker rm -f #{container} &>/dev/null || :)
    end

    (
      get_containers.call("jancajthaml/rest-contract-test-ramltestee")
    ).flatten.each { |container| teardown_container.call(container) }

    Dir.glob("/var/log/contract_*").each { |file|
      puts "deleting file #{file}"
      #FileUtils.mkdir_p folder
      FileUtils.rm_rf file
    }

    print "[ suite ended    ]"
  end

end
