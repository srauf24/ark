import { Outlet } from "react-router-dom";
import { SignOutButton } from "@clerk/clerk-react";
import { LogOut } from "lucide-react";
import { Button } from "@/components/ui/button";

export function Layout() {
    return (
        <div className="min-h-screen bg-background font-sans antialiased">
            <header className="sticky top-0 z-50 w-full border-b border-border-subtle bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
                <div className="container flex h-14 items-center">
                    <div className="mr-4 hidden md:flex">
                        <a className="mr-6 flex items-center space-x-2" href="/">
                            <img src="/icon-192.png" alt="Ark Logo" className="h-6 w-6" />
                            <span className="hidden font-bold sm:inline-block">Ark</span>
                        </a>
                    </div>
                    <div className="ml-auto flex items-center space-x-4">
                        <SignOutButton>
                            <Button variant="ghost" size="sm" className="text-muted-foreground hover:text-foreground">
                                <LogOut className="mr-2 h-4 w-4" />
                                Sign Out
                            </Button>
                        </SignOutButton>
                    </div>
                </div>
            </header>
            <main className="container py-6">
                <Outlet />
            </main>
        </div>
    );
}
