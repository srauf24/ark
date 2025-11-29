import { Link } from "react-router-dom";
import {
  Server,
  Activity,
  LayoutDashboard,
  Zap,
  Shield,
  Network,
  ChevronRight,
} from "lucide-react";
import { Button } from "@/components/ui/button";

export function LandingPage() {
  return (
    <div className="min-h-screen bg-[#0a0a0a] text-white overflow-hidden">
      {/* Background grid pattern */}
      <div className="fixed inset-0 opacity-20 pointer-events-none">
        <div
          className="absolute inset-0"
          style={{
            backgroundImage: `
              linear-gradient(to right, #1a1a1a 1px, transparent 1px),
              linear-gradient(to bottom, #1a1a1a 1px, transparent 1px)
            `,
            backgroundSize: "40px 40px",
          }}
        />
      </div>

      {/* Animated scan line */}
      <div className="fixed inset-0 pointer-events-none">
        <div className="scan-line" />
      </div>

      {/* Navigation */}
      <nav className="relative z-10 border-b border-[#1a1a1a] bg-black/50 backdrop-blur-md">
        <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gradient-to-br from-emerald-400 to-cyan-500 rounded-md flex items-center justify-center font-mono font-bold text-black text-lg shadow-lg shadow-emerald-500/50">
              ARK
            </div>
            <span className="font-bold text-xl tracking-tight">
              ARK
            </span>
          </div>
          <div className="flex items-center gap-4">
            <Link to="/sign-in">
              <Button
                variant="outline"
                className="border-emerald-500/30 text-emerald-400 hover:bg-emerald-500/10 hover:border-emerald-500/50"
              >
                Sign In
              </Button>
            </Link>
          </div>
        </div>
      </nav>

      <main className="relative z-10">
        {/* Hero Section */}
        <section className="max-w-7xl mx-auto px-6 pt-24 pb-32">
          <div className="grid lg:grid-cols-2 gap-16 items-center">
            {/* Left: Content */}
            <div className="space-y-8">
              {/* Status indicator */}
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/30 rounded-full">
                  <div className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse-glow" />
                  <span className="text-emerald-400 text-xs font-mono uppercase tracking-wider">
                    System Online
                  </span>
                </div>
              </div>

              <div className="space-y-6">
                <h1 className="text-6xl lg:text-7xl font-black leading-tight tracking-tight">
                  Track Your
                  <br />
                  <span className="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 via-cyan-400 to-blue-400 animate-gradient">
                    Homelab
                  </span>
                  <br />
                  Infrastructure
                </h1>
                <p className="text-xl text-gray-400 leading-relaxed max-w-xl">
                  Manage your servers, VMs, and network equipment. Document every
                  change. Search your infrastructure history. Simple and powerful.
                </p>
              </div>

              <div className="flex flex-col sm:flex-row gap-4">
                <Link to="/sign-up">
                  <Button
                    size="lg"
                    className="bg-gradient-to-r from-emerald-500 to-cyan-500 hover:from-emerald-600 hover:to-cyan-600 text-black font-bold shadow-lg shadow-emerald-500/30 group"
                  >
                    Get Started
                    <ChevronRight className="ml-2 w-4 h-4 group-hover:translate-x-1 transition-transform" />
                  </Button>
                </Link>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-2 gap-8 pt-8 max-w-sm">
                {[
                  { value: "Unlimited", label: "Assets" },
                  { value: "Free", label: "Full Access" },
                ].map((stat, i) => (
                  <div
                    key={i}
                    className="group hover:scale-105 transition-transform"
                  >
                    <div className="text-2xl font-bold font-mono text-emerald-400 group-hover:text-cyan-400 transition-colors">
                      {stat.value}
                    </div>
                    <div className="text-sm text-gray-500 uppercase tracking-wider">
                      {stat.label}
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Right: Browser mockup showing asset detail page */}
            <div className="relative">
              <div className="relative bg-[#0d0d0d] border border-gray-800 rounded-lg overflow-hidden shadow-2xl shadow-emerald-500/10">
                {/* Browser header */}
                <div className="bg-[#1a1a1a] border-b border-gray-800 px-4 py-3 flex items-center gap-2">
                  <div className="flex gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-red-500/80" />
                    <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
                    <div className="w-3 h-3 rounded-full bg-emerald-500/80" />
                  </div>
                  <div className="flex-1 ml-3">
                    <div className="bg-[#0a0a0a] border border-gray-700 rounded px-3 py-1 text-xs text-gray-500 max-w-md">
                      localhost:3000/assets/abc-123
                    </div>
                  </div>
                </div>

                {/* Browser content - Asset detail page */}
                <div className="p-6 space-y-4 max-h-[400px] overflow-hidden">
                  {/* Asset header */}
                  <div className="flex items-start justify-between animate-fade-in">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 rounded-lg bg-emerald-500/10 border border-emerald-500/30 flex items-center justify-center">
                        <Server className="w-6 h-6 text-emerald-400" />
                      </div>
                      <div>
                        <h2 className="text-lg font-bold text-white">prod-server-01</h2>
                        <div className="flex items-center gap-2 mt-1">
                          <span className="px-2 py-0.5 bg-emerald-500/10 border border-emerald-500/30 rounded text-emerald-400 text-xs">
                            Server
                          </span>
                          <span className="text-xs text-gray-500">
                            prod-server-01.home.lab
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Asset metadata */}
                  <div className="bg-black border border-gray-800 rounded-lg p-4 animate-fade-in" style={{ animationDelay: "200ms" }}>
                    <div className="text-xs text-gray-400 mb-2 font-mono">METADATA</div>
                    <div className="space-y-1 text-sm font-mono">
                      <div className="flex gap-2">
                        <span className="text-cyan-400">cpu:</span>
                        <span className="text-gray-300">Intel Xeon E5-2670</span>
                      </div>
                      <div className="flex gap-2">
                        <span className="text-cyan-400">ram:</span>
                        <span className="text-gray-300">64GB DDR4</span>
                      </div>
                      <div className="flex gap-2">
                        <span className="text-cyan-400">os:</span>
                        <span className="text-gray-300">Ubuntu 24.04 LTS</span>
                      </div>
                    </div>
                  </div>

                  {/* Activity Log section */}
                  <div className="animate-fade-in" style={{ animationDelay: "400ms" }}>
                    <div className="flex items-center justify-between mb-3">
                      <h3 className="text-sm font-bold text-white flex items-center gap-2">
                        <Activity className="w-4 h-4 text-emerald-400" />
                        Activity Logs
                      </h3>
                      <button className="px-3 py-1 bg-emerald-500/10 border border-emerald-500/30 rounded text-emerald-400 text-xs font-medium hover:bg-emerald-500/20 transition-colors">
                        Add Log
                      </button>
                    </div>

                    {/* Log entry */}
                    <div className="bg-black border border-gray-800 rounded-lg p-4 space-y-3">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="text-sm text-gray-300 leading-relaxed">
                            Updated nginx configuration to enable gzip compression and
                            HTTP/2. Restarted service successfully.
                          </div>
                          <div className="flex flex-wrap gap-2 mt-3">
                            <span className="px-2 py-0.5 bg-cyan-500/10 border border-cyan-500/30 rounded text-cyan-400 text-xs">
                              nginx
                            </span>
                            <span className="px-2 py-0.5 bg-cyan-500/10 border border-cyan-500/30 rounded text-cyan-400 text-xs">
                              performance
                            </span>
                            <span className="px-2 py-0.5 bg-cyan-500/10 border border-cyan-500/30 rounded text-cyan-400 text-xs">
                              config
                            </span>
                          </div>
                        </div>
                      </div>
                      <div className="flex items-center gap-3 pt-2 border-t border-gray-800">
                        <span className="text-xs text-gray-500">2 hours ago</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Decorative elements */}
              <div className="absolute -top-12 -right-12 w-64 h-64 bg-emerald-500/20 rounded-full blur-3xl animate-pulse-slow" />
              <div className="absolute -bottom-12 -left-12 w-64 h-64 bg-cyan-500/20 rounded-full blur-3xl animate-pulse-slow" style={{ animationDelay: "1s" }} />
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section className="max-w-7xl mx-auto px-6 py-24">
          <div className="text-center space-y-4 mb-16">
            <div className="inline-flex items-center gap-2 px-3 py-1.5 bg-cyan-500/10 border border-cyan-500/30 rounded-full">
              <Zap className="w-4 h-4 text-cyan-400" />
              <span className="text-cyan-400 text-xs font-mono uppercase tracking-wider">
                Features
              </span>
            </div>
            <h2 className="text-4xl lg:text-5xl font-black">
              Everything You Need
            </h2>
            <p className="text-gray-400 text-lg max-w-2xl mx-auto">
              Purpose-built for managing complex homelab environments with
              precision and style.
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
            {[
              {
                icon: Server,
                title: "Infrastructure Inventory",
                description:
                  "Track servers, VMs, containers, network devices, and storage. Store custom metadata like CPU, RAM, IP addresses, and specs for each asset.",
                color: "emerald",
              },
              {
                icon: Activity,
                title: "Configuration Logs",
                description:
                  "Document every config change, upgrade, and troubleshooting session. Add tags for easy categorization. Searchable history for all your infrastructure.",
                color: "cyan",
              },
              {
                icon: LayoutDashboard,
                title: "Clean Web Interface",
                description:
                  "Responsive dashboard with forms, modals, and real-time updates. Works seamlessly on desktop, tablet, and mobile devices.",
                color: "blue",
              },
              {
                icon: Shield,
                title: "Simple & Secure",
                description:
                  "Sign up and start tracking your infrastructure immediately. No installation required. Secure authentication with your data safely stored.",
                color: "purple",
              },
            ].map((feature, i) => (
              <div
                key={i}
                className="group relative bg-gradient-to-b from-[#111] to-black border border-gray-800 rounded-lg p-6 hover:border-gray-700 transition-all duration-300 hover:scale-105 hover:shadow-xl hover:shadow-emerald-500/10"
              >
                <div className={`w-12 h-12 rounded-lg bg-${feature.color}-500/10 border border-${feature.color}-500/30 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform`}>
                  <feature.icon className={`w-6 h-6 text-${feature.color}-400`} />
                </div>
                <h3 className="text-xl font-bold mb-2 group-hover:text-emerald-400 transition-colors">
                  {feature.title}
                </h3>
                <p className="text-gray-500 leading-relaxed">
                  {feature.description}
                </p>

                {/* Hover glow effect */}
                <div className="absolute inset-0 rounded-lg bg-gradient-to-br from-emerald-500/0 to-cyan-500/0 group-hover:from-emerald-500/5 group-hover:to-cyan-500/5 transition-all duration-300 pointer-events-none" />
              </div>
            ))}
          </div>
        </section>

        {/* Trust Section */}
        <section className="max-w-7xl mx-auto px-6 py-24">
          <div className="bg-gradient-to-br from-[#111] to-black border border-gray-800 rounded-2xl p-12 lg:p-16 relative overflow-hidden">
            {/* Background decorations */}
            <div className="absolute top-0 right-0 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl" />
            <div className="absolute bottom-0 left-0 w-96 h-96 bg-cyan-500/10 rounded-full blur-3xl" />

            <div className="relative z-10 max-w-4xl mx-auto">
              <div className="text-center space-y-6 mb-12">
                <div className="inline-flex items-center gap-2 px-3 py-1.5 bg-emerald-500/10 border border-emerald-500/30 rounded-full">
                  <Shield className="w-4 h-4 text-emerald-400" />
                  <span className="text-emerald-400 text-xs font-mono uppercase tracking-wider">
                    Built for Homelabbers
                  </span>
                </div>
                <h2 className="text-4xl lg:text-5xl font-black">
                  Start Tracking
                  <br />
                  <span className="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 to-cyan-400">
                    In Minutes
                  </span>
                </h2>
                <p className="text-gray-400 text-lg max-w-2xl mx-auto">
                  No installation, no configuration. Just sign up and start
                  documenting your homelab infrastructure.
                </p>
              </div>

              <div className="grid md:grid-cols-2 gap-12 max-w-2xl mx-auto">
                {[
                  {
                    icon: Zap,
                    title: "Instant Access",
                    description: "Sign up and start using ARK immediately. No downloads, no setup required.",
                  },
                  {
                    icon: Network,
                    title: "Homelab Native",
                    description: "Designed for the complexities of real homelab setups. Built by homelabbers, for homelabbers.",
                  },
                ].map((item, i) => (
                  <div
                    key={i}
                    className="text-center space-y-3 group hover:scale-105 transition-transform"
                  >
                    <div className="w-16 h-16 mx-auto rounded-xl bg-gradient-to-br from-emerald-500/10 to-cyan-500/10 border border-gray-700 flex items-center justify-center group-hover:border-emerald-500/50 transition-colors">
                      <item.icon className="w-8 h-8 text-emerald-400 group-hover:text-cyan-400 transition-colors" />
                    </div>
                    <h3 className="text-xl font-bold">{item.title}</h3>
                    <p className="text-gray-500 leading-relaxed">
                      {item.description}
                    </p>
                  </div>
                ))}
              </div>

              <div className="mt-12 text-center">
                <Link to="/sign-up">
                  <Button
                    size="lg"
                    className="bg-gradient-to-r from-emerald-500 to-cyan-500 hover:from-emerald-600 hover:to-cyan-600 text-black font-bold shadow-lg shadow-emerald-500/30"
                  >
                    Get Started Free
                    <ChevronRight className="ml-2 w-4 h-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </section>

        {/* Footer */}
        <footer className="border-t border-gray-900 bg-black/50 backdrop-blur-md">
          <div className="max-w-7xl mx-auto px-6 py-12">
            <div className="flex flex-col md:flex-row justify-between items-start gap-8">
              <div className="space-y-4">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 bg-gradient-to-br from-emerald-400 to-cyan-500 rounded-md flex items-center justify-center font-mono font-bold text-black text-lg shadow-lg shadow-emerald-500/50">
                    ARK
                  </div>
                  <span className="font-bold text-xl">ARK</span>
                </div>
                <p className="text-gray-500 text-sm leading-relaxed max-w-xs">
                  Homelab infrastructure management, built with pride.
                </p>
              </div>

              <div className="flex items-center gap-2 text-sm">
                <div className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse-glow" />
                <span className="text-gray-500 font-mono">
                  All systems operational
                </span>
              </div>
            </div>

            <div className="pt-8 mt-8 border-t border-gray-900">
              <p className="text-gray-500 text-sm font-mono text-center">
                © 2025 ARK. All rights reserved.
              </p>
            </div>
          </div>
        </footer>
      </main>
    </div>
  );
}
