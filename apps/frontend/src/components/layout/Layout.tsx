import { Outlet } from "react-router-dom";

export function Layout() {
    return (
        <div className="min-h-screen bg-background font-sans antialiased">
            <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
                <div className="container flex h-14 items-center">
                    <div className="mr-4 hidden md:flex">
                        <a className="mr-6 flex items-center space-x-2" href="/">
                            <img src="/icon-192.png" alt="Ark Logo" className="h-6 w-6" />
                            <span className="hidden font-bold sm:inline-block">Ark</span>
                        </a>
                    </div>
                </div>
            </header>
            <main className="container py-6">
                <Outlet />
            </main>
        </div>
    );
}
