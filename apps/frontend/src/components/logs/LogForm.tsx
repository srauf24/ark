import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Button } from "@/components/ui/button";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { TagBadge } from "@/components/common/TagBadge";
import { useState } from "react";
import { ZCreateLogRequest } from "@ark/zod";

type CreateLogRequest = z.infer<typeof ZCreateLogRequest>;

interface LogFormProps {
    defaultValues?: CreateLogRequest;
    onSubmit: (data: CreateLogRequest) => void;
    isSubmitting: boolean;
}

export function LogForm({ defaultValues, onSubmit, isSubmitting }: LogFormProps) {
    const form = useForm<CreateLogRequest>({
        resolver: zodResolver(ZCreateLogRequest),
        defaultValues: defaultValues || {
            content: "",
            tags: [],
        },
    });

    const [tagInput, setTagInput] = useState("");

    const handleTagKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter" || e.key === ",") {
            e.preventDefault();
            const newTag = tagInput.trim();
            if (newTag) {
                const currentTags = form.getValues("tags") || [];
                if (!currentTags.includes(newTag)) {
                    form.setValue("tags", [...currentTags, newTag]);
                }
                setTagInput("");
            }
        }
    };

    const removeTag = (tagToRemove: string) => {
        const currentTags = form.getValues("tags") || [];
        form.setValue(
            "tags",
            currentTags.filter((tag) => tag !== tagToRemove)
        );
    };

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                    control={form.control}
                    name="content"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Content</FormLabel>
                            <FormControl>
                                <Textarea
                                    placeholder="Describe the change or event..."
                                    className="min-h-[120px] resize-y"
                                    {...field}
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                <FormField
                    control={form.control}
                    name="tags"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>Tags</FormLabel>
                            <FormControl>
                                <div className="space-y-2">
                                    <Input
                                        placeholder="Type tag and press Enter..."
                                        value={tagInput}
                                        onChange={(e) => setTagInput(e.target.value)}
                                        onKeyDown={handleTagKeyDown}
                                    />
                                    <div className="flex flex-wrap gap-2">
                                        {field.value?.map((tag) => (
                                            <TagBadge
                                                key={tag}
                                                tag={tag}
                                                onRemove={() => removeTag(tag)}
                                            />
                                        ))}
                                    </div>
                                </div>
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                <div className="flex justify-end gap-2">
                    <Button type="submit" disabled={isSubmitting}>
                        {isSubmitting ? "Saving..." : "Save Log"}
                    </Button>
                </div>
            </form>
        </Form>
    );
}
