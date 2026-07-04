namespace MyRestApi.Models
{
    public class TournTbrk
    {
        public int TournId { get; set; }
        public string Tourn {get; set; }
    }

    public class TournToFrom
    {
        public int TournIdFrom { get; set; }
        public int TournIdTo {get; set; }
    }
}