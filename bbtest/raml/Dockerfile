FROM library/ruby:latest

RUN gem install sinatra require_all

COPY ./app.rb /opt/app/entrypoint.rb

ENTRYPOINT ["ruby", "/opt/app/entrypoint.rb"]
CMD ["-p", "8080", "-o", "0.0.0.0"]
