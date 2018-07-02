
step "contract is run with following parameres" do |parameters|
  args = parameters.split("\n").map { |line|
    line.strip().split(" ")
  }.flatten.reject(&:empty?).join(" ")

  $call_id += 1
  %x(contract #{args} >/var/log/contract_#{$call_id}.log 2>&1)
  expect($?).to be_success, %x(cat /var/log/contract_#{$call_id}.log)
end
