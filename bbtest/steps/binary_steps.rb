
step "contract is run with following parameres" do |parameters|
  args = parameters.split("\n").map { |line|
    line.strip().split(" ")
  }.flatten.reject(&:empty?).join(" ")

  $call_id += 1
  std = %x(contract #{args} >/var/log/contract_#{$call_id}.log 2>&1)
  expect($?).to be_success, std
end
