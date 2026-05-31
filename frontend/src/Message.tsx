import { useQuery } from "@tanstack/react-query";

export default function Message() {
    const query = useQuery({
        queryKey: ["message"],
        queryFn: async () => {
            const response = await fetch("/api/v1/message");
            return response.json();
        },
    })

    if (query.isFetching) {
        return <p>Loading...</p>
    }

    if (query.error) {
        return <p>Failed to fetch message.</p>
    }

    return <p>{query.data.message}</p>
}
