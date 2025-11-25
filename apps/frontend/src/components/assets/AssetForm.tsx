import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ZCreateAssetRequest, ZUpdateAssetRequest } from "@ark/zod";
import { Loader2 } from "lucide-react";
import {
    Form,
    FormControl,
    FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import type { Asset } from "@/types";

interface AssetFormProps {
    asset?: Asset;
    onSuccess: () => void;
    onCancel: () => void;
    isPending?: boolean;
}

const ASSET_TYPES = [
    { value: "server", label: "Server" },
    { value: "vm", label: "Virtual Machine" },
    { value: "nas", label: "NAS" },
    { value: "container", label: "Container" },
    { value: "network", label: "Network Device" },
    { value: "other", label: "Other" },
] as const;

export function AssetForm({
    asset,
    onSuccess,
    onCancel,
    isPending = false,
}: AssetFormProps) {
    const mode = asset ? "edit" : "create";
    const schema = mode === "edit" ? ZUpdateAssetRequest : ZCreateAssetRequest;

    const form = useForm<{
        name: string;
        type?: string;
        hostname?: string;
        metadata?: string;
    }>({
        resolver: zodResolver(schema),
        defaultValues: {
            name: asset?.name || "",
            type: asset?.type || undefined,
            hostname: asset?.hostname || "",
            metadata: asset?.metadata ? JSON.stringify(asset.metadata, null, 2) : "",
        },
    });

    const onSubmit = (data: any) => {
        // Parse metadata JSON if provided
        let parsedData = { ...data };
        if (data.metadata && data.metadata.trim()) {
            try {
                parsedData.metadata = JSON.parse(data.metadata);
            } catch (error) {
                form.setError("metadata", {
                    type: "manual",
                    message: "Invalid JSON format",
                });
                return;
            }
        } else {
            // Remove empty metadata field
            delete parsedData.metadata;
        }

        // Remove empty hostname field
        if (!parsedData.hostname || parsedData.hostname.trim() === "") {
            delete parsedData.hostname;
        }

        // This will be connected to API mutations in Step 6
        console.log("Form submitted:", parsedData);
        onSuccess();
    };

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                {/* Name Field */}
                <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Name *</FormLabel>
                            <FormControl>
                                <Input
                                    placeholder="e.g., Production Server 1"
                                    {...field}
                                    disabled={isPending}
                                />
                            </FormControl>
                            <FormDescription>
                                A unique name for this asset (1-100 characters)
                            </FormDescription>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                {/* Type Field */}
                <FormField
                    control={form.control}
                    name="type"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Type</FormLabel>
                            <Select
                                onValueChange={field.onChange}
                                defaultValue={field.value}
                                disabled={isPending}
                            >
                                <FormControl>
                                    <SelectTrigger>
                                        <SelectValue placeholder="Select asset type" />
                                    </SelectTrigger>
                                </FormControl>
                                <SelectContent>
                                    {ASSET_TYPES.map((type) => (
                                        <SelectItem key={type.value} value={type.value}>
                                            {type.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                            <FormDescription>
                                The category of this asset (optional)
                            </FormDescription>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                {/* Hostname Field */}
                <FormField
                    control={form.control}
                    name="hostname"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Hostname</FormLabel>
                            <FormControl>
                                <Input
                                    placeholder="e.g., prod-server-01.example.com"
                                    {...field}
                                    disabled={isPending}
                                />
                            </FormControl>
                            <FormDescription>
                                Network hostname or FQDN (optional, max 255 characters)
                            </FormDescription>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                {/* Metadata Field */}
                <FormField
                    control={form.control}
                    name="metadata"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Metadata (JSON)</FormLabel>
                            <FormControl>
                                <Textarea
                                    placeholder='{"cpu": "Intel Xeon", "ram": "64GB"}'
                                    className="font-mono text-sm"
                                    rows={6}
                                    {...field}
                                    disabled={isPending}
                                />
                            </FormControl>
                            <FormDescription>
                                Additional information in JSON format (optional)
                            </FormDescription>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                {/* Form Actions */}
                <div className="flex justify-end gap-3">
                    <Button
                        type="button"
                        variant="outline"
                        onClick={onCancel}
                        disabled={isPending}
                    >
                        Cancel
                    </Button>
                    <Button type="submit" disabled={isPending}>
                        {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        {mode === "create" ? "Create Asset" : "Save Changes"}
                    </Button>
                </div>
            </form>
        </Form>
    );
}
