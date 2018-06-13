
step "contract is run with following parameres" do |parameters|
  # fixme time this
  puts %x(contract #{parameters.strip})
end
