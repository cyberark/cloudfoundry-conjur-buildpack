require 'roda'

class TestApp < Roda
  plugin :default_headers, 'Content-Type'=>'text/html'

  route do |r|
    r.root do
      "
      <h1>Visit us @ www.conjur.org!</h1>

      <h3>Space-wide Secrets</h3>
      <p>Database Username: #{ENV['SPACE_USERNAME']}</p>
      <p>Database Password: #{ENV['SPACE_PASSWORD']}</p>
      "
    end
  end
end
