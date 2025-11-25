import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { LogForm } from "./LogForm";
import { useCreateLog, useUpdateLog } from "@/hooks/useLogs";
import { ZAssetLog } from "@ark/zod";
import { z } from "zod";

type AssetLog = z.infer<typeof ZAssetLog>;

interface LogFormModalProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    initialData?: AssetLog;
    assetId: string;
}

export function LogFormModal({
    open,
    onOpenChange,
    initialData,
    assetId,
}: LogFormModalProps) {
    const createLog = useCreateLog(assetId);
    const updateLog = useUpdateLog(assetId);

    const isEditing = !!initialData;
    const isSubmitting = createLog.isPending || updateLog.isPending;

    const handleSubmit = async (data: { content: string; tags?: string[] }) => {
        try {
            if (isEditing && initialData) {
                await updateLog.mutateAsync({
                    id: initialData.id,
                    data,
                });
            } else {
                await createLog.mutateAsync(data);
            }
            onOpenChange(false);
        } catch (error) {
            // Error handling is done in the hook via toast
            console.error("Failed to save log:", error);
        }
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>{isEditing ? "Edit Log" : "Add Log"}</DialogTitle>
                </DialogHeader>
                <LogForm
                    defaultValues={
                        initialData
                            ? {
                                content: initialData.content,
                                tags: initialData.tags || [],
                            }
                            : undefined
                    }
                    onSubmit={handleSubmit}
                    isSubmitting={isSubmitting}
                />
            </DialogContent>
        </Dialog>
    );
}
