import { type Asset } from "@ark/zod";
import { Link } from "react-router-dom";
import { format } from "date-fns";
import {
    Server,
    Monitor,
    HardDrive,
    Box,
    Network,
    HelpCircle,
    type LucideIcon
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface AssetCardProps {
    asset: Asset;
}

const getAssetIcon = (type: string | null | undefined): LucideIcon => {
    switch (type) {
        case "server":
            return Server;
        case "vm":
            return Monitor;
        case "nas":
            return HardDrive;
        case "container":
            return Box;
        case "network":
            return Network;
        default:
            return HelpCircle;
    }
};

export function AssetCard({ asset }: AssetCardProps) {
    const Icon = getAssetIcon(asset.type);

    return (
        <Link to={`/assets/${asset.id}`} className="block transition-transform hover:scale-[1.02]">
            <Card className="h-full cursor-pointer hover:border-primary/50 hover:bg-gradient-to-br hover:from-surface-1 hover:to-surface-2 transition-colors duration-300">
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium truncate pr-2">
                        {asset.name}
                    </CardTitle>
                    <Icon className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                    <div className="flex flex-col gap-2">
                        {asset.hostname && (
                            <div className="text-xs text-muted-foreground truncate">
                                {asset.hostname}
                            </div>
                        )}
                        <div className="flex items-center justify-between mt-2">
                            <Badge variant="secondary" className="text-xs font-normal">
                                {asset.type || "other"}
                            </Badge>
                            <span className="text-[10px] text-muted-foreground">
                                {format(new Date(asset.updated_at), "MMM d, yyyy")}
                            </span>
                        </div>
                    </div>
                </CardContent>
            </Card>
        </Link>
    );
}
