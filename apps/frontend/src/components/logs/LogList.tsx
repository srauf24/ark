import { useState } from "react"; // Trigger HMR
import { Plus, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { LogEntry } from "./LogEntry";
import { LogFormModal } from "./LogFormModal";
import { DeleteConfirmDialog } from "@/components/common/DeleteConfirmDialog";
import { useLogs, useDeleteLog } from "@/hooks/useLogs";
import { ZAssetLog } from "@ark/zod";
import { z } from "zod";

type AssetLog = z.infer<typeof ZAssetLog>;

interface LogListProps {
    assetId: string;
}

export function LogList({ assetId }: LogListProps) {
    const { data, isLoading, isError, refetch } = useLogs(assetId);
    const deleteLog = useDeleteLog(assetId);

    const [isFormOpen, setIsFormOpen] = useState(false);
    const [editingLog, setEditingLog] = useState<AssetLog | undefined>(undefined);
    const [deletingLogId, setDeletingLogId] = useState<string | null>(null);

    const handleCreate = () => {
        setEditingLog(undefined);
        setIsFormOpen(true);
    };

    const handleEdit = (log: AssetLog) => {
        setEditingLog(log);
        setIsFormOpen(true);
    };

    const handleDelete = (id: string) => {
        setDeletingLogId(id);
    };

    const confirmDelete = async () => {
        if (deletingLogId) {
            await deleteLog.mutateAsync(deletingLogId);
            setDeletingLogId(null);
        }
    };

    if (isLoading) {
        return (
            <div className="space-y-4">
                <div className="flex items-center justify-between">
                    <h2 className="text-lg font-semibold">Logs</h2>
                    <Skeleton className="h-9 w-24" />
                </div>
                <div className="space-y-4">
                    {[1, 2, 3].map((i) => (
                        <Skeleton key={i} className="h-24 w-full" />
                    ))}
                </div>
            </div>
        );
    }

    if (isError) {
        return (
            <div className="space-y-4">
                <div className="flex items-center justify-between">
                    <h2 className="text-lg font-semibold">Logs</h2>
                </div>
                <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertTitle>Error</AlertTitle>
                    <AlertDescription className="flex items-center gap-2">
                        Failed to load logs.
                        <Button variant="link" className="p-0 h-auto font-normal text-destructive underline" onClick={() => refetch()}>
                            Try again
                        </Button>
                    </AlertDescription>
                </Alert>
            </div>
        );
    }

    const logs = data?.logs || [];

    return (
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <h2 className="text-lg font-semibold">Logs</h2>
                <Button onClick={handleCreate} size="sm" className="gap-1">
                    <Plus className="h-4 w-4" />
                    Add Log
                </Button>
            </div>

            {logs.length === 0 ? (
                <div className="rounded-lg border border-dashed p-8 text-center">
                    <p className="text-muted-foreground mb-2">No logs recorded yet.</p>
                    <Button variant="outline" size="sm" onClick={handleCreate}>
                        Create your first log
                    </Button>
                </div>
            ) : (
                <div className="space-y-4">
                    {logs.map((log) => (
                        <LogEntry
                            key={log.id}
                            log={log}
                            onEdit={() => handleEdit(log)}
                            onDelete={() => handleDelete(log.id)}
                        />
                    ))}
                </div>
            )}

            <LogFormModal
                open={isFormOpen}
                onOpenChange={setIsFormOpen}
                initialData={editingLog}
                assetId={assetId}
            />

            <DeleteConfirmDialog
                isOpen={!!deletingLogId}
                onCancel={() => setDeletingLogId(null)}
                onConfirm={confirmDelete}
                title="Delete Log"
                message="Are you sure you want to delete this log entry? This action cannot be undone."
                isLoading={deleteLog.isPending}
            />
        </div>
    );
}
