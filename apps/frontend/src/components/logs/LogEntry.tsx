import { formatDistanceToNow } from "date-fns";
import { Edit2, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { TagBadge } from "@/components/common/TagBadge";
import { ZAssetLog } from "@ark/zod";
import { z } from "zod";

type AssetLog = z.infer<typeof ZAssetLog>;

interface LogEntryProps {
    log: AssetLog;
    onEdit: () => void;
    onDelete: () => void;
}

export function LogEntry({ log, onEdit, onDelete }: LogEntryProps) {
    return (
        <div className="group relative flex flex-col gap-2 rounded-lg border p-4 hover:bg-muted/50 transition-colors">
            <div className="flex items-start justify-between gap-4">
                <div className="flex-1 space-y-1">
                    <p className="whitespace-pre-wrap text-sm leading-relaxed text-foreground/90">
                        {log.content}
                    </p>

                    {log.tags && log.tags.length > 0 && (
                        <div className="flex flex-wrap gap-2 pt-2">
                            {log.tags.map((tag) => (
                                <TagBadge key={tag} tag={tag} />
                            ))}
                        </div>
                    )}

                    <p className="text-xs text-muted-foreground pt-1">
                        {formatDistanceToNow(new Date(log.created_at), { addSuffix: true })}
                    </p>
                </div>

                <div className="flex items-center gap-1 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity">
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-muted-foreground hover:text-foreground"
                        onClick={onEdit}
                        aria-label="Edit log"
                    >
                        <Edit2 className="h-4 w-4" />
                    </Button>
                    <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-muted-foreground hover:text-destructive"
                        onClick={onDelete}
                        aria-label="Delete log"
                    >
                        <Trash2 className="h-4 w-4" />
                    </Button>
                </div>
            </div>
        </div>
    );
}
