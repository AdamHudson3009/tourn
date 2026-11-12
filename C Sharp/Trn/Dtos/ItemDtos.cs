namespace Trn.Dtos;

// JSON input
public record ItemRequest(int TrnId, List<int> ConfId, List<List<int>> PlyrsConIds);


// JSON output
public record MatchResponse(List<int> Plyrs1Ids, List<int> Plyrs2Ids);

// String output since database entry (or not)
public record OutResponse(string msg);