
step "logs contains following" do |args|
  lines = args.split("\n").map(&:strip).reject(&:empty?)

  abspath = "/var/log/contract_#{$call_id}.log"

  unless File.file?(abspath)
    raise "file:  #{abspath} is a directory" if File.directory?(abspath)
    raise "file:  #{abspath} was not found\nfiles: #{Dir[File.dirname(abspath)+"/*"]}"
  end

  contents = File.open(abspath, 'rb').read.split("\n").map(&:strip).reject(&:empty?)

  puts "expected lines: #{lines}"
  puts "actual logs: #{contents}"
  #puts contents

  #lines.

  #{}%x(docker logs #{container} >/reports/#{label}.log 2>&1)

  #containers = %x(docker ps -a -f -f name=#{label} | awk '{ print $1,$2 }' | sed 1,1d)



  #std = %x(contract #{args})
  #expect($?).to be_success, std
end
