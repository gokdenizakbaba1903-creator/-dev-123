"use client";

import { useState, useEffect } from 'react';
import { Search, ShieldAlert, CheckCircle, Bug, Lock, Server, Globe, Activity, History, Trash2, Download, Info, ChevronRight, AlertTriangle } from 'lucide-react';
import { Radar, RadarChart, PolarGrid, PolarAngleAxis, PolarRadiusAxis, ResponsiveContainer } from 'recharts';

export default function Home() {
  const [url, setUrl] = useState('');
  const [scanning, setScanning] = useState(false);
  const [result, setResult] = useState<any>(null);
  const [history, setHistory] = useState<any[]>([]);
  const [expandedFinding, setExpandedFinding] = useState<number | null>(null);

  useEffect(() => {
    const savedHistory = localStorage.getItem('secscan_history');
    if (savedHistory) setHistory(JSON.parse(savedHistory));
  }, []);

  const saveToHistory = (newResult: any) => {
    const updatedHistory = [newResult, ...history.filter(h => h.id !== newResult.id)].slice(0, 10);
    setHistory(updatedHistory);
    localStorage.setItem('secscan_history', JSON.stringify(updatedHistory));
  };

  const startScan = async () => {
    if (!url) return;
    setScanning(true);
    setResult(null);
    setExpandedFinding(null);
    try {
      const res = await fetch('http://localhost:8080/api/scan', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      });
      const data = await res.json();
      if (data.id) connectSSE(data.id);
    } catch (err) {
      console.error(err);
      setScanning(false);
    }
  };

  const connectSSE = (id: string) => {
    const eventSource = new EventSource(`http://localhost:8080/stream/${id}`);
    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setResult(data);
      if (data.status === 'Completed') {
        eventSource.close();
        setScanning(false);
        saveToHistory(data);
      }
    };
  };

  const clearHistory = () => {
    setHistory([]);
    localStorage.removeItem('secscan_history');
  };

  const radarData = result?.radar_data ? Object.entries(result.radar_data).map(([key, value]) => ({
    subject: key,
    A: value,
    fullMark: 100,
  })) : [
    { subject: 'Port', A: 0 }, { subject: 'Headers', A: 0 }, { subject: 'TLS', A: 0 },
    { subject: 'XSS', A: 0 }, { subject: 'SQLi', A: 0 }, { subject: 'CVE', A: 0 },
  ];

  return (
    <div className="flex h-screen bg-[#0a0c10] text-gray-100 overflow-hidden">
      {/* SIDEBAR */}
      <aside className="w-72 bg-white/5 border-r border-white/10 flex flex-col p-6 hidden lg:flex">
        <div className="flex items-center gap-3 mb-10">
          <div className="p-2 bg-blue-600 rounded-lg">
            <ShieldAlert size={24} />
          </div>
          <span className="text-xl font-bold tracking-tight">SecScan <span className="text-blue-500">AI</span></span>
        </div>

        <div className="flex-1 overflow-y-auto space-y-4">
          <div className="flex justify-between items-center text-xs font-semibold text-gray-500 uppercase tracking-widest px-2">
            <span>Recent Scans</span>
            <button onClick={clearHistory} className="hover:text-red-400 transition-colors"><Trash2 size={14}/></button>
          </div>
          {history.map((h, i) => (
            <button 
              key={i} 
              onClick={() => setResult(h)}
              className="w-full text-left p-3 rounded-xl hover:bg-white/5 border border-transparent hover:border-white/10 transition-all group"
            >
              <div className="text-sm font-medium truncate mb-1 group-hover:text-blue-400">{h.url}</div>
              <div className="flex justify-between items-center text-[10px] text-gray-500">
                <span>{h.score} Score</span>
                <span>{h.findings.length} findings</span>
              </div>
            </button>
          ))}
          {history.length === 0 && <div className="text-center py-10 text-gray-600 text-sm italic">No history yet</div>}
        </div>

        <div className="pt-6 border-t border-white/10 text-xs text-gray-500">
          v2.1 Premium Enterprise
        </div>
      </aside>

      {/* MAIN CONTENT */}
      <main className="flex-1 overflow-y-auto p-8 relative">
        <div className="max-w-5xl mx-auto">
          <header className="flex justify-between items-end mb-10">
            <div>
              <h1 className="text-4xl font-bold mb-2">Security Dashboard</h1>
              <p className="text-gray-500">Real-time vulnerability assessment and reporting.</p>
            </div>
            <div className="flex gap-4">
              <div className="glass-card flex items-center gap-3 px-4 py-2 text-sm">
                <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
                <span>System Status: Online</span>
              </div>
            </div>
          </header>

          <section className="glass-card p-4 flex gap-4 items-center mb-10 border-blue-500/20 bg-blue-500/5">
            <div className="relative flex-1">
              <Globe className="absolute left-4 top-1/2 -translate-y-1/2 text-blue-500/60" size={20} />
              <input 
                type="text" 
                placeholder="https://example.com"
                className="w-full bg-transparent border-none py-4 pl-12 pr-4 focus:outline-none text-lg font-medium"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
              />
            </div>
            <button 
              onClick={startScan}
              disabled={scanning}
              className="bg-blue-600 hover:bg-blue-700 hover:shadow-[0_0_20px_rgba(37,99,235,0.4)] px-8 py-4 rounded-xl font-bold flex items-center gap-3 transition-all disabled:opacity-50 active:scale-95"
            >
              {scanning ? <Activity className="animate-spin" /> : <Search size={22} />}
              {scanning ? 'Analyzing Target...' : 'Start Scan'}
            </button>
          </section>

          {(scanning || result) && (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-12">
              <div className="lg:col-span-2 space-y-8">
                {/* PROGRESS CARD */}
                <div className="glass-card p-6 overflow-hidden relative group">
                  {scanning && <div className="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-blue-600 to-emerald-400 animate-shimmer" />}
                  <div className="flex justify-between items-center mb-6">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-blue-500/10 rounded-lg text-blue-500">
                        <Activity size={20} />
                      </div>
                      <h2 className="font-bold text-lg">Scan Intelligence</h2>
                    </div>
                    <div className="flex flex-col items-end">
                      <span className="text-3xl font-black text-blue-400">{result?.progress || 0}%</span>
                      <span className="text-[10px] text-gray-500 uppercase tracking-widest">{result?.status}</span>
                    </div>
                  </div>
                  <div className="w-full bg-white/5 h-3 rounded-full mb-2 p-0.5">
                    <div 
                      className="bg-gradient-to-r from-blue-600 to-blue-400 h-full rounded-full transition-all duration-1000 shadow-[0_0_10px_rgba(59,130,246,0.5)]"
                      style={{ width: `${result?.progress || 0}%` }}
                    />
                  </div>
                </div>

                {/* FINDINGS CARD */}
                <div className="glass-card p-6 min-h-[400px]">
                  <h2 className="text-xl font-bold mb-6 flex items-center gap-2">
                    <ShieldAlert className="text-blue-500" /> Detected Vulnerabilities 
                    <span className="ml-2 px-2 py-0.5 bg-white/5 rounded text-xs font-normal text-gray-400">{result?.findings?.length || 0} found</span>
                  </h2>
                  <div className="space-y-4">
                    {result?.findings?.map((f: any, i: number) => (
                      <div 
                        key={i} 
                        className={`group border rounded-2xl transition-all overflow-hidden ${expandedFinding === i ? 'bg-white/10 border-white/20' : 'bg-white/[0.02] border-white/5 hover:border-white/10'}`}
                      >
                        <div 
                          className="flex items-center gap-4 p-4 cursor-pointer"
                          onClick={() => setExpandedFinding(expandedFinding === i ? null : i)}
                        >
                          <div className={`p-3 rounded-xl ${f.severity === 'High' ? 'bg-red-500/10 text-red-500 shadow-[inset_0_0_10px_rgba(239,68,68,0.2)]' : 'bg-emerald-500/10 text-emerald-500 shadow-[inset_0_0_10px_rgba(16,185,129,0.2)]'}`}>
                            {f.severity === 'High' ? <Bug size={20} /> : <Lock size={20} />}
                          </div>
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <span className="font-bold text-gray-200">{f.module}</span>
                              <span className={`text-[10px] px-2 py-0.5 rounded-full font-bold uppercase ${f.severity === 'High' ? 'bg-red-500/20 text-red-400' : 'bg-emerald-500/20 text-emerald-400'}`}>
                                {f.severity}
                              </span>
                            </div>
                            <div className="text-sm text-gray-500 line-clamp-1">{f.message}</div>
                          </div>
                          <ChevronRight className={`text-gray-600 transition-transform ${expandedFinding === i ? 'rotate-90' : ''}`} />
                        </div>
                        
                        {expandedFinding === i && (
                          <div className="px-4 pb-4 pt-2 border-t border-white/5 animate-fade-in">
                            <div className="bg-black/20 rounded-xl p-4 space-y-4">
                              <div>
                                <div className="text-xs font-bold text-gray-400 uppercase tracking-wide mb-1 flex items-center gap-1">
                                  <Info size={12}/> Observation
                                </div>
                                <p className="text-sm text-gray-300">{f.message}</p>
                              </div>
                              <div className="p-3 bg-blue-500/5 rounded-lg border border-blue-500/10">
                                <div className="text-xs font-bold text-blue-400 uppercase tracking-wide mb-1 flex items-center gap-1">
                                  <CheckCircle size={12}/> Remediation Strategy
                                </div>
                                <p className="text-sm text-blue-100/70">
                                  {f.severity === 'High' ? 'Implement strict input validation, sanitize all contextual output, and use parameterized queries. Review relevant security patches immediately.' : 'Enhance security configuration, enforce strict headers, and periodically audit the components involved.'}
                                </p>
                              </div>
                            </div>
                          </div>
                        )}
                      </div>
                    ))}
                    {(!result?.findings || result.findings.length === 0) && !scanning && (
                      <div className="text-center py-20 bg-emerald-500/5 rounded-3xl border border-emerald-500/10">
                        <CheckCircle className="mx-auto text-emerald-500 mb-2" size={40} />
                        <p className="font-bold text-emerald-400">Target Appears Secure</p>
                        <p className="text-sm text-emerald-900/60">No high-risk vulnerabilities detected in current modules.</p>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              <div className="space-y-8">
                {/* FINAL SCORE CARD */}
                <div className="glass-card p-8 text-center relative group overflow-hidden">
                  <div className="absolute inset-0 bg-gradient-to-b from-blue-500/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
                  <h2 className="text-gray-500 mb-4 uppercase text-[10px] font-bold tracking-[0.2em]">Compliance Score</h2>
                  <div className={`text-8xl font-black mb-6 transition-all duration-500 ${result?.score === 'F' ? 'text-red-500 drop-shadow-[0_0_20px_rgba(239,68,68,0.5)]' : 'text-emerald-400 drop-shadow-[0_0_20px_rgba(52,211,153,0.5)]'}`}>
                    {result?.score || '--'}
                  </div>
                  {result?.status === 'Completed' && (
                    <button 
                      onClick={() => window.open(`http://localhost:8080/api/report/${result.id}/pdf`, '_blank')}
                      className="w-full bg-white/10 hover:bg-blue-600 text-white py-4 rounded-2xl font-bold transition-all flex items-center justify-center gap-2 border border-white/10 hover:border-blue-500 shadow-xl active:scale-95"
                    >
                      <Download size={18} />
                      Export PDF Report
                    </button>
                  )}
                </div>

                {/* RADAR CARD */}
                <div className="glass-card p-6 h-[380px] flex flex-col">
                  <h2 className="text-center font-bold text-gray-400 text-sm tracking-wide mb-4">Risk Distribution</h2>
                  <div className="flex-1 w-full">
                    <ResponsiveContainer width="100%" height="100%">
                      <RadarChart cx="50%" cy="50%" outerRadius="75%" data={radarData}>
                        <PolarGrid stroke="#333" />
                        <PolarAngleAxis dataKey="subject" tick={{ fill: '#666', fontSize: 10, fontWeight: 'bold' }} />
                        <Radar
                          name="Security"
                          dataKey="A"
                          stroke="#2563eb"
                          fill="#3b82f6"
                          fillOpacity={0.6}
                        />
                      </RadarChart>
                    </ResponsiveContainer>
                  </div>
                </div>

                {/* SYSTEM INFO */}
                <div className="glass-card p-6 text-sm space-y-4">
                   <h3 className="font-bold text-gray-300 flex items-center gap-2"><Globe size={16} className="text-blue-500"/> Target Intel</h3>
                   <div className="space-y-2 text-xs">
                     <div className="flex justify-between border-b border-white/5 pb-2">
                       <span className="text-gray-500">Scan Origin</span>
                       <span className="text-gray-300">Ireland (AWS-EU)</span>
                     </div>
                     <div className="flex justify-between border-b border-white/5 pb-2">
                       <span className="text-gray-500">Engine Version</span>
                       <span className="text-gray-300">v4.0.2-Enterprise</span>
                     </div>
                     <div className="flex justify-between">
                       <span className="text-gray-500">Detection Logic</span>
                       <span className="text-gray-300">AI-Enhanced Pattern</span>
                     </div>
                   </div>
                </div>
              </div>
            </div>
          )}

          {!scanning && !result && (
             <div className="flex flex-col items-center justify-center py-32 text-center opacity-40">
                <div className="w-24 h-24 bg-blue-500/10 rounded-full flex items-center justify-center mb-6">
                  <ShieldAlert size={48} className="text-blue-500" />
                </div>
                <h2 className="text-3xl font-bold mb-2">Ready to Scan</h2>
                <p className="max-w-md">Input a target URL above to begin a deep perimeter security assessment.</p>
             </div>
          )}
        </div>
      </main>
    </div>
  );
}
