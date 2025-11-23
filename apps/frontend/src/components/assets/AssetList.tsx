import { useApiClient } from "@/api";
import { useQuery } from "@tanstack/react-query";
import { AssetCard } from "./AssetCard";
import { Loader2, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import { Link } from "react-router-dom";

export function AssetList() {
    const apiClient = useApiClient();

    const { data, isLoading, isError } = useQuery({
        queryKey: ["assets"],
        queryFn: async () => {
            const response = await apiClient.Assets.listAssets({
                query: { limit: 100 },
            });

            if (response.status !== 200) {
                throw new Error("Failed to fetch assets");
            }

            return response.body;
        },
    });

    if (isLoading) {
        return (
            <div className="flex h-[50vh] w-full items-center justify-center" role="status">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    if (isError) {
        return (
            <Alert variant="destructive">
                <AlertCircle className="h-4 w-4" />
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>
                    Failed to load assets. Please try again later.
                </AlertDescription>
            </Alert>
        );
    }

    if (!data?.assets || data.assets.length === 0) {
        return (
            <div className="flex h-[50vh] flex-col items-center justify-center gap-4 text-center">
                <div className="text-lg font-semibold">No assets found</div>
                <p className="text-sm text-muted-foreground">
                    Get started by adding your first asset to the inventory.
                </p>
                <Button asChild>
                    <Link to="/assets/new">
                        <Plus className="mr-2 h-4 w-4" />
                        Add Asset
                    </Link>
                </Button>
            </div>
        );
    }

    return (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {data.assets.map((asset) => (
                <AssetCard key={asset.id} asset={asset} />
            ))}
        </div>
    );
}
