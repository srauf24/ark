import { Badge } from "@/components/ui/badge";
import { X } from "lucide-react";
import { cn } from "@/lib/utils";

interface TagBadgeProps {
    tag: string;
    onRemove?: () => void;
    className?: string;
}

export function TagBadge({ tag, onRemove, className }: TagBadgeProps) {
    return (
        <Badge
            variant="secondary"
            className={cn(
                "px-2 py-0.5 text-xs font-medium flex items-center gap-1",
                className
            )}
        >
            {tag}
            {onRemove && (
                <button
                    type="button"
                    onClick={(e) => {
                        e.stopPropagation();
                        onRemove();
                    }}
                    className="hover:bg-muted-foreground/20 rounded-full p-0.5 transition-colors"
                    aria-label={`Remove ${tag} tag`}
                >
                    <X className="h-3 w-3" />
                </button>
            )}
        </Badge>
    );
}
