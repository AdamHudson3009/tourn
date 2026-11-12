using Dapper;
using Trn.Data;
using Trn.Dtos;
using Microsoft.Extensions.Configuration;

namespace Trn.Services;

public static class WCPrelim
{
    private static readonly DbConnectionFactory _db = new DbConnectionFactory(
        new ConfigurationBuilder()
            .SetBasePath(AppContext.BaseDirectory)
            .AddJsonFile("appsettings.json", optional: false)
            .Build());
    
    public static (string, List<int>) ProcessWCPrelim(int trnId, List<int> plyrsIds, int wcBase, int maxRnd, int rowRnd) {
        MatchResponse mtchResponse = wc_preliminary(trnId, plyrsIds, wcBase);

        if (rowRnd <= maxRnd) {
            if (!WCCommon.wc_insert(trnId, maxRnd+1, mtchResponse.Plyrs1Ids, mtchResponse.Plyrs2Ids)) {
                return ("Prelim: Failed to insert matches.", new List<int>());
            }
            return ("Prelim: Inserted into prelims", new List<int>());
        }
        
        (bool isValid, List<int> wcWns) = WCCommon.wc_validate(trnId, maxRnd+1, mtchResponse.Plyrs1Ids, mtchResponse.Plyrs2Ids);
        if (!isValid) {
            return ("Prelim: Matches not consistent with flow.", new List<int>());
        }
            
        for (int indx = 0; indx < wcWns.Count; indx++) {
            if (wcWns[indx] < 0)
                return ("Prelim: Matches not completed.", new List<int>());
        }

        mtchResponse.Plyrs1Ids.Clear();
        mtchResponse.Plyrs2Ids.Clear();
        return ("", wcWns);
    }

    private static MatchResponse wc_preliminary(int trnId, List<int> plyrsIds, int wcBase) {
        var plyrTmp = new List<int>();

        int wcBaseT2 = wcBase*2;
        int PlyrsIdsValid = plyrsIds.Count;
        if (PlyrsIdsValid > wcBaseT2) PlyrsIdsValid = wcBaseT2;

        int start = wcBaseT2-PlyrsIdsValid;
        int endIndx = PlyrsIdsValid - 1;
        for (int indx = start; indx <= endIndx; indx++) {
            plyrTmp.Add(plyrsIds[indx]);
        }

        return WCCommon.wc_round_setup(plyrTmp);
    }
}