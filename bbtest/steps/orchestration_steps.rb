step "no :container :label is running" do |container, label|

  containers = %x(docker ps -a -f name=#{label} | awk '{ print $1,$2 }' | grep #{container} | awk '{print $1 }' 2>/dev/null)
  expect($?).to be_success

  ids = containers.split("\n").map(&:strip).reject(&:empty?)

  return if ids.empty?

  ids.each { |id|
    eventually(timeout: 3) {
      puts "wanting to kill #{id}"
      send ":container running state is :state", id, false

      label = %x(docker inspect --format='{{.Name}}' #{id})
      label = ($? == 0 ? label.strip : id)

      %x(docker logs #{id} >/reports#{label}.log 2>&1)
      %x(docker rm -f #{id} &>/dev/null || :)
    }
  }
end

step ":container running state is :state" do |container, state|
  eventually(timeout: 5) {
    %x(docker #{state ? "start" : "stop"} #{container} >/dev/null 2>&1)
    container_state = %x(docker inspect -f {{.State.Running}} #{container} 2>/dev/null)
    expect($?).to be_success
    expect(container_state.strip).to eq(state ? "true" : "false")
  }
end

step ":container :version is started with" do |container, version, label, params|
  containers = %x(docker ps -a -f status=running -f name=#{label} | awk '{ print $1,$2 }' | sed 1,1d)
  expect($?).to be_success
  containers = containers.split("\n").map(&:strip).reject(&:empty?)

  unless containers.empty?
    id, image = containers[0].split(" ")
    return if (image == "#{container}:#{version}")
  end

  send "no :container :label is running", container, label

  prefix = ENV.fetch('COMPOSE_PROJECT_NAME', "")
  my_id = %x(cat /etc/hostname).strip
  args = [
    "docker",
    "run",
    "-d",
    "--net=#{prefix}_default",
    "--volumes-from=#{my_id}",
    "--log-driver=json-file",
    "-h #{label}",
    "--net-alias=#{label}",
    "--name=#{label}"
  ] << params << [
    "#{container}:#{version}",
    "2>&1"
  ]

  id = %x(#{args.join(" ")}).strip
  expect($?).to be_success, id

  eventually(timeout: 10) {
    send ":container running state is :state", id, true
  }

  eventually(timeout: 10) {
    logs = %x(docker logs #{id} 2>&1 | grep "WEBrick::HTTPServer#start: pid=")
    expect($?).to be_success
    expect(logs.split("\n").map(&:strip).reject(&:empty?)).not_to be_empty
  }
end

step ":label is not running" do |label|
  containers = %x(docker ps -a -f status=running -f name=#{label} | awk '{ print $1 }' | sed 1,1d)
  expect($?).to be_success
  containers = containers.split("\n").map(&:strip).reject(&:empty?)

  return if containers.empty?

  id = containers[0]

  eventually(timeout: 10) {
    send ":container running state is :state", id, false
  }
end

step ":label is running" do |label|
  image = "jancajthaml/rest-contract-test-#{label}"
  send ":container :version is started with", image, ENV.fetch("VERSION", "latest"), label, [
    "-p 8080"
  ]
end
