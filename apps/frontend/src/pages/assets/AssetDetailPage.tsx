import { useState } from "react";
import { useApiClient } from "@/api";
import { useQuery } from "@tanstack/react-query";
import { useParams, Link, useNavigate } from "react-router-dom";
import { format } from "date-fns";
import {
    ArrowLeft,
    Server,
    Monitor,
    HardDrive,
    Box,
    Network,
    HelpCircle,
    Calendar,
    Clock,
    Edit,
    Trash2,
    type LucideIcon,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { AlertCircle, Loader2 } from "lucide-react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { AssetForm } from "@/components/assets/AssetForm";
import { DeleteConfirmDialog } from "@/components/common/DeleteConfirmDialog";
import { useDeleteAsset } from "@/hooks/useAssets";

const getAssetIcon = (type: string | null | undefined): LucideIcon => {
    switch (type) {
        case "server": return Server;
        case "vm": return Monitor;
        case "nas": return HardDrive;
        case "container": return Box;
        case "network": return Network;
        default: return HelpCircle;
    }
};

export function AssetDetailPage() {
    const { id } = useParams<{ id: string }>();
    const apiClient = useApiClient();
    const navigate = useNavigate();
    const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

    const deleteMutation = useDeleteAsset();

    const { data, isLoading, isError } = useQuery({
        queryKey: ["asset", id],
        queryFn: async () => {
            if (!id) throw new Error("Asset ID is required");
            const response = await apiClient.Assets.getAssetById({
                params: { id },
            });

            if (response.status !== 200) {
                throw new Error("Failed to fetch asset");
            }

            return response.body.data;
        },
        enabled: !!id,
    });

    const handleDelete = () => {
        if (!id) return;

        deleteMutation.mutate(id, {
            onSuccess: () => {
                navigate("/assets");
            },
        });
    };

    if (isLoading) {
        return (
            <div className="flex h-[50vh] w-full items-center justify-center" role="status">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    if (isError || !data) {
        return (
            <div className="space-y-4">
                <Button variant="ghost" asChild className="pl-0">
                    <Link to="/assets">
                        <ArrowLeft className="mr-2 h-4 w-4" />
                        Back to Assets
                    </Link>
                </Button>
                <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertTitle>Error</AlertTitle>
                    <AlertDescription>
                        Failed to load asset details. The asset may not exist or you don't
                        have permission to view it.
                    </AlertDescription>
                </Alert>
            </div>
        );
    }

    const Icon = getAssetIcon(data.type);

    return (
        <>
            <div className="space-y-6">
                <div className="flex items-center justify-between">
                    <div className="space-y-1">
                        <Button variant="ghost" asChild className="pl-0 -ml-2 mb-2">
                            <Link to="/assets">
                                <ArrowLeft className="mr-2 h-4 w-4" />
                                Back to Assets
                            </Link>
                        </Button>
                        <div className="flex items-center gap-3">
                            <div className="flex h-10 w-10 items-center justify-center rounded-lg border bg-card text-card-foreground shadow-sm">
                                <Icon className="h-6 w-6 text-muted-foreground" />
                            </div>
                            <div>
                                <h1 className="text-2xl font-bold tracking-tight">{data.name}</h1>
                                {data.hostname && (
                                    <p className="text-sm text-muted-foreground">{data.hostname}</p>
                                )}
                            </div>
                        </div>
                    </div>
                    <div className="flex gap-2">
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={() => setIsEditDialogOpen(true)}
                        >
                            <Edit className="mr-2 h-4 w-4" />
                            Edit
                        </Button>
                        <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => setIsDeleteDialogOpen(true)}
                        >
                            <Trash2 className="mr-2 h-4 w-4" />
                            Delete
                        </Button>
                    </div>
                </div>

                <div className="grid gap-6 md:grid-cols-2">
                    <Card>
                        <CardHeader>
                            <CardTitle>Details</CardTitle>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div className="space-y-1">
                                    <p className="text-muted-foreground">Type</p>
                                    <Badge variant="secondary">{data.type || "other"}</Badge>
                                </div>
                                <div className="space-y-1">
                                    <p className="text-muted-foreground">ID</p>
                                    <p className="font-mono text-xs">{data.id}</p>
                                </div>
                                <div className="space-y-1">
                                    <p className="text-muted-foreground flex items-center gap-1">
                                        <Calendar className="h-3 w-3" /> Created
                                    </p>
                                    <p>{format(new Date(data.created_at), "PPP p")}</p>
                                </div>
                                <div className="space-y-1">
                                    <p className="text-muted-foreground flex items-center gap-1">
                                        <Clock className="h-3 w-3" /> Updated
                                    </p>
                                    <p>{format(new Date(data.updated_at), "PPP p")}</p>
                                </div>
                            </div>
                        </CardContent>
                    </Card>

                    <Card>
                        <CardHeader>
                            <CardTitle>Metadata</CardTitle>
                        </CardHeader>
                        <CardContent>
                            {data.metadata && Object.keys(data.metadata).length > 0 ? (
                                <pre className="overflow-auto rounded-lg bg-muted p-4 text-xs font-mono">
                                    {JSON.stringify(data.metadata, null, 2)}
                                </pre>
                            ) : (
                                <p className="text-sm text-muted-foreground italic">
                                    No metadata available
                                </p>
                            )}
                        </CardContent>
                    </Card>
                </div>

                {/* Placeholder for Logs - to be implemented in ARK-F4 */}
                <Card className="opacity-60">
                    <CardHeader>
                        <CardTitle>Logs</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex h-32 items-center justify-center rounded-md border border-dashed">
                            <p className="text-sm text-muted-foreground">
                                Log history will be available soon
                            </p>
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Edit Dialog */}
            <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
                <DialogContent className="max-w-2xl">
                    <DialogHeader>
                        <DialogTitle>Edit Asset</DialogTitle>
                        <DialogDescription>
                            Update the details of this asset.
                        </DialogDescription>
                    </DialogHeader>
                    <AssetForm
                        asset={data}
                        onSuccess={() => setIsEditDialogOpen(false)}
                        onCancel={() => setIsEditDialogOpen(false)}
                    />
                </DialogContent>
            </Dialog>

            {/* Delete Confirmation Dialog */}
            <DeleteConfirmDialog
                isOpen={isDeleteDialogOpen}
                onConfirm={handleDelete}
                onCancel={() => setIsDeleteDialogOpen(false)}
                title="Delete Asset"
                message={`Are you sure you want to delete "${data.name}"? This action cannot be undone.`}
                isLoading={deleteMutation.isPending}
            />
        </>
    );
}
