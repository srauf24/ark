import { useState } from "react";
import { AssetCard } from "./AssetCard";
import { AssetForm } from "./AssetForm";
import { Loader2, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle } from "lucide-react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";

import { useAssets } from "@/hooks/useAssets";

export function AssetList() {
    const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);

    const { data, isLoading, isError } = useAssets({ limit: 100 });

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
            <>
                <div className="flex h-[50vh] flex-col items-center justify-center gap-4 text-center">
                    <div className="text-lg font-semibold">No assets found</div>
                    <p className="text-sm text-muted-foreground">
                        Get started by adding your first asset to the inventory.
                    </p>
                    <Button onClick={() => setIsCreateDialogOpen(true)}>
                        <Plus className="mr-2 h-4 w-4" />
                        Add Asset
                    </Button>
                </div>

                <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                    <DialogContent className="max-w-2xl">
                        <DialogHeader>
                            <DialogTitle>Create New Asset</DialogTitle>
                            <DialogDescription>
                                Add a new asset to your inventory. Fill in the details below.
                            </DialogDescription>
                        </DialogHeader>
                        <AssetForm
                            onSuccess={() => setIsCreateDialogOpen(false)}
                            onCancel={() => setIsCreateDialogOpen(false)}
                        />
                    </DialogContent>
                </Dialog>
            </>
        );
    }

    return (
        <>
            <div className="mb-6 flex items-center justify-between">
                <div>
                    <h2 className="text-2xl font-bold tracking-tight">Assets</h2>
                    <p className="text-sm text-muted-foreground">
                        Manage your infrastructure assets
                    </p>
                </div>
                <Button onClick={() => setIsCreateDialogOpen(true)}>
                    <Plus className="mr-2 h-4 w-4" />
                    Add Asset
                </Button>
            </div>

            <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {data.assets.map((asset) => (
                    <AssetCard key={asset.id} asset={asset} />
                ))}
            </div>

            <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                <DialogContent className="max-w-2xl">
                    <DialogHeader>
                        <DialogTitle>Create New Asset</DialogTitle>
                        <DialogDescription>
                            Add a new asset to your inventory. Fill in the details below.
                        </DialogDescription>
                    </DialogHeader>
                    <AssetForm
                        onSuccess={() => setIsCreateDialogOpen(false)}
                        onCancel={() => setIsCreateDialogOpen(false)}
                    />
                </DialogContent>
            </Dialog>
        </>
    );
}
