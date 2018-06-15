require 'sinatra'

set :show_exceptions, true

require 'pathname'
require 'json'
require 'date'

################################################################################
# configuration
################################################################################

#helpers Sinatra::Param

configure do
  set :raise_errors => true
  set :logging, true
end

################################################################################
# hooks
################################################################################
before do
  #content_type :json
  logger.level = 0
  @req = nil
  if request.body.size > 0
    request.body.rewind
    @req = JSON.parse(request.body.read)
  end
end

################################################################################
# state
################################################################################

$g_people = Hash.new

################################################################################
# api
################################################################################
get '/ping' do
  #logger.debug "hit GET /ping"

  [ 200, { "time" => DateTime.now.strftime('%Q') }.to_json ]
end

get '/:version/person' do

  # fixme list all make it optional
  return [ 400, {} ] unless params.has_key?("LastName")

  filterRule = params["LastName"]

  people = $g_people.values.select { |v|
    v[:LastName] == filterRule
  }.map { |v|
    { "firstName" => v[:FirstName], "lastName" => v[:LastName] }
  }

  [ 200, people.to_json ]
end

post '/:version/person' do
  #logger.debug "hit POST /{version}/person"

  return [ 400, {} ] if @req.nil?

  logger.debug "plaintext playload #{@req}"

  first, last, *rest = @req["name"].split(/ /)
  id = "#{first}/#{last}".to_i(36).to_s

  logger.debug "new person #{@req["name"]} id: #{id}"

  [ 419, {} ] if $g_people.key?(id)

  $g_people[id] = {
    :FirstName => first,
    :LastName => last
  }

  resp = { "id" => id }

  return [ 200, resp.to_json ]
end

get '/:version/person/:id' do
  #logger.debug "hit GET /{version}/person/{id}"

  id = params["id"]

  [ 417, {} ] unless $g_people.key?(id)

  person = $g_people[id]
  resp = { "firstName" => person[:FirstName], "lastName" => person[:LastName] }

  logger.debug "person #{id} found"

  return [ 200, resp.to_json ]
end

delete '/:version/person/:id' do
  #logger.debug "hit DELETE /{version}/person/{id}"

  id = params["id"]

  [ 417, {} ] unless $g_people.key?(id)

  $g_people.delete(id)

  return [ 200, {}.to_json ]
end
