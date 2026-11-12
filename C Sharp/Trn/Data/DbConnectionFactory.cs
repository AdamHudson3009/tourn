using MySqlConnector;

namespace Trn.Data;

public class DbConnectionFactory
{
    private readonly string _connectionString;

    public DbConnectionFactory(IConfiguration config) =>
        _connectionString = config.GetConnectionString("Default")!;

    public MySqlConnection Create() => new MySqlConnection(_connectionString);
}