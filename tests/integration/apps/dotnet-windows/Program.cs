var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/", () => {
    var body = $"""
    <h1>Visit us @ www.conjur.org!</h1>
    <h3>Space-wide Secrets</h3>
    <p>Database Username: {app.Configuration["SPACE_USERNAME"]}</p>
    <p>Database Password: {app.Configuration["SPACE_PASSWORD"]}</p>
    """;

    return body;
});

app.Run();
