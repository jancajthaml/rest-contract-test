
step "logs contains following" do |args|
  lines = args.split("\n").map(&:strip).reject(&:empty?)

  abspath = "/var/log/contract_#{$call_id}.log"

  unless File.file?(abspath)
    raise "file:  #{abspath} is a directory" if File.directory?(abspath)
    raise "file:  #{abspath} was not found\nfiles: #{Dir[File.dirname(abspath)+"/*"]}"
  end

  contents = File.open(abspath, 'rb').read.split("\n").map(&:strip).reject(&:empty?)
  lines.each { |line|
    found = false
    contents.each { |l|
      next unless l.include? line
      found = true
      break
    }
    raise "#{line} was not found in logs:\n#{contents}" unless found
  }
end
