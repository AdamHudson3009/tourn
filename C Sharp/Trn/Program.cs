using Trn.Data;
using Trn.Dtos;
using Trn.Services;

var builder = WebApplication.CreateBuilder(args);

// Dependency injection
builder.Services.AddSingleton<DbConnectionFactory>();
builder.Services.AddScoped<WCService>();

// CORS: allow the Angular dev server
builder.Services.AddCors(options =>
{
    options.AddPolicy("ng", policy =>
        policy.SetIsOriginAllowed(origin => new Uri(origin).Host == "localhost")
            .AllowAnyHeader()
            .AllowAnyMethod());
});

var app = builder.Build();

app.UseCors("ng");   // must come before the endpoints are mapped

app.MapPost("/trn/wc", async (ItemRequest req, WCService svc) =>
    Results.Ok(await svc.ProcessWCAsync(req)));

app.Run();