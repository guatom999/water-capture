import { useEffect, useState, useMemo } from "react"
import { Link, useSearchParams } from 'react-router-dom';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend, ReferenceLine, Brush } from 'recharts';
import Footer from "./Footer";
import { getMapMarkerDetailService } from "../services/waterLevelService";
import type { LocationDetail, WaterDetailResponse } from "../types/waterDetail";

const SectionDetail = () => {
    const [searchParams] = useSearchParams();
    const stationId = searchParams.get('station_id');

    const [sectionDetail, setSectionDetail] = useState<WaterDetailResponse | null>(null);
    const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
    const [dateRange, setDateRange] = useState<'all' | '1day' | '7days' | '30days'>('1day');

    const getSectionDetail = async () => {
        if (stationId) {
            const response = await getMapMarkerDetailService(stationId);
            setSectionDetail(response as WaterDetailResponse)
        }
    };

    useEffect(() => {
        getSectionDetail();
    }, [stationId]);

    // Get bank level from API
    const bankLevel = sectionDetail?.markers?.bank_level ?? 1.29;

    // Filter data by date range
    const filteredMarkers = useMemo(() => {
        if (!sectionDetail?.markers?.detail) return [];

        const now = new Date();

        return sectionDetail.markers.detail.filter((marker: LocationDetail) => {
            if (!marker.measured_at) return false;
            const markerDate = new Date(marker.measured_at);

            switch (dateRange) {
                case '1day':
                    const oneDayAgo = new Date(now.getTime() - 1 * 24 * 60 * 60 * 1000);
                    return markerDate >= oneDayAgo;
                case '7days':
                    const sevenDaysAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
                    return markerDate >= sevenDaysAgo;
                case '30days':
                    const thirtyDaysAgo = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
                    return markerDate >= thirtyDaysAgo;
                default:
                    return true;
            }
        });
    }, [sectionDetail?.markers?.detail, dateRange]);

    // Calculate statistics
    const statistics = useMemo(() => {
        if (filteredMarkers.length === 0) return null;

        const levels = filteredMarkers
            .map((m: LocationDetail) => m.level_cm)
            .filter((l): l is number => l !== null);

        if (levels.length === 0) return null;

        return {
            max: Math.max(...levels),
            min: Math.min(...levels),
            avg: levels.reduce((a: number, b: number) => a + b, 0) / levels.length,
            count: levels.length
        };
    }, [filteredMarkers]);

    useEffect(() => {
        if (sectionDetail) {
            console.log('Markers data:', sectionDetail.markers);
        }
    }, [sectionDetail]);

    // Generate Y-axis ticks: 0, 0.5, 1.0... up to bankLevel, then bankLevel, then .0/.5 up to bankLevel + 1
    const generateYAxisTicks = (bankLevelValue: number): number[] => {
        const ticks: number[] = [];
        const maxValue = bankLevelValue + 1;

        // Add ticks at 0.5 intervals from 0 up to (but not exceeding) bankLevel
        for (let i = 0; i <= maxValue; i += 0.5) {
            const rounded = Math.round(i * 10) / 10;
            if (rounded < bankLevelValue) {
                ticks.push(rounded);
            }
        }

        // Add the exact bankLevel value
        ticks.push(bankLevelValue);

        // Add .0 and .5 ticks above bankLevel up to bankLevel + 1
        const nextHalf = Math.ceil(bankLevelValue * 2) / 2;
        for (let i = nextHalf; i <= maxValue; i += 0.5) {
            const rounded = Math.round(i * 10) / 10;
            if (rounded > bankLevelValue && !ticks.includes(rounded)) {
                ticks.push(rounded);
            }
        }

        return ticks.sort((a, b) => a - b);
    };

    const getDangerColor = (danger: string | null) => {
        switch (danger?.toUpperCase()) {
            case 'CRITICAL':
                return 'bg-red-100 text-red-800 border-red-300';
            case 'DANGER':
                return 'bg-red-100 text-red-800 border-red-300';
            case 'WATCH':
                return 'bg-yellow-100 text-yellow-800 border-yellow-300';
            case 'SAFE':
                return 'bg-green-100 text-green-800 border-green-300';
            default:
                return 'bg-gray-100 text-gray-800 border-gray-300';
        }
    };

    // Format just the time (HH:mm)
    const formatTime = (dateString: string | null) => {
        if (!dateString) return 'N/A';
        return new Date(dateString).toLocaleString('th-TH', {
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    // Format just the date for section header
    const formatDateHeader = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('th-TH', {
            weekday: 'long',
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    };

    // Get date key for grouping (YYYY-MM-DD)
    const getDateKey = (dateString: string | null) => {
        if (!dateString) return 'unknown';
        return new Date(dateString).toISOString().split('T')[0];
    };

    // Group markers by date
    const groupByDate = (markers: LocationDetail[]) => {
        const sorted = [...markers].sort((a, b) => {
            const dateA = a.measured_at ? new Date(a.measured_at).getTime() : 0;
            const dateB = b.measured_at ? new Date(b.measured_at).getTime() : 0;
            return sortOrder === 'desc' ? dateB - dateA : dateA - dateB;
        });

        const groups: { [key: string]: LocationDetail[] } = {};
        sorted.forEach(marker => {
            const key = getDateKey(marker.measured_at);
            if (!groups[key]) groups[key] = [];
            groups[key].push(marker);
        });

        return Object.entries(groups);
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-cyan-50">
            {/* Header */}
            <div className="bg-blue-600 text-white py-8 shadow-lg">
                <Link to="/">
                    <div className="container mx-auto px-4">
                        <h1 className="text-3xl font-bold mb-2">Water Level Details</h1>
                        <p className="text-blue-100">Station ID: {stationId}</p>
                    </div>
                </Link>

            </div>

            {sectionDetail ? (
                <div className="container mx-auto px-4 py-8">
                    {sectionDetail.markers.detail.length === 0 ? (
                        <div className="bg-white rounded-lg shadow-md p-8 text-center">
                            <div className="text-gray-400 text-6xl mb-4"></div>
                            <p className="text-gray-600 text-lg">No water level data found for this location.</p>
                        </div>
                    ) : (
                        <div className="space-y-6">
                            {/* Date Range Filter */}
                            <div className="bg-white rounded-xl shadow-lg p-4 border border-gray-100">
                                <div className="flex flex-wrap items-center gap-4">
                                    <span className="text-sm font-medium text-gray-600">ช่วงเวลา:</span>
                                    <div className="flex gap-2">
                                        {[
                                            { key: 'all', label: 'ทั้งหมด' },
                                            { key: '1day', label: '24 ชม.' },
                                            { key: '7days', label: '7 วัน' },
                                            { key: '30days', label: '30 วัน' }
                                        ].map(option => (
                                            <button
                                                key={option.key}
                                                onClick={() => setDateRange(option.key as typeof dateRange)}
                                                className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${dateRange === option.key
                                                    ? 'bg-blue-500 text-white'
                                                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                                                    }`}
                                            >
                                                {option.label}
                                            </button>
                                        ))}
                                    </div>
                                </div>
                            </div>

                            {/* Statistics Cards */}
                            {statistics && (
                                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                                    <div className="bg-white rounded-xl shadow p-4 border border-gray-100">
                                        <div className="text-sm text-gray-500">ค่าสูงสุด</div>
                                        <div className="text-2xl font-bold text-red-500">{statistics.max.toFixed(1)} MSL</div>
                                    </div>
                                    <div className="bg-white rounded-xl shadow p-4 border border-gray-100">
                                        <div className="text-sm text-gray-500">ค่าต่ำสุด</div>
                                        <div className="text-2xl font-bold text-green-500">{statistics.min.toFixed(1)} MSL</div>
                                    </div>
                                    <div className="bg-white rounded-xl shadow p-4 border border-gray-100">
                                        <div className="text-sm text-gray-500">ค่าเฉลี่ย</div>
                                        <div className="text-2xl font-bold text-blue-500">{statistics.avg.toFixed(1)} MSL</div>
                                    </div>
                                    <div className="bg-white rounded-xl shadow p-4 border border-gray-100">
                                        <div className="text-sm text-gray-500">จำนวนข้อมูล</div>
                                        <div className="text-2xl font-bold text-gray-700">{statistics.count} รายการ</div>
                                    </div>
                                </div>
                            )}

                            {/* Water Level Chart */}
                            <div className="bg-white rounded-xl shadow-lg p-6 border border-gray-100">
                                <h3 className="text-lg font-semibold text-gray-800 mb-4">Water Level Over Time</h3>
                                <div className="h-96">
                                    <ResponsiveContainer width="100%" height="100%">
                                        <LineChart
                                            data={[...filteredMarkers]
                                                .sort((a, b) => {
                                                    const dateA = a.measured_at ? new Date(a.measured_at).getTime() : 0;
                                                    const dateB = b.measured_at ? new Date(b.measured_at).getTime() : 0;
                                                    return dateA - dateB;
                                                })
                                                .map((m) => ({
                                                    time: m.measured_at ? new Date(m.measured_at).toLocaleString('th-TH', {
                                                        hour: '2-digit',
                                                        minute: '2-digit'
                                                    }) : 'N/A',
                                                    fullDate: m.measured_at ? new Date(m.measured_at).toLocaleString('th-TH', {
                                                        weekday: 'short',
                                                        year: 'numeric',
                                                        month: 'short',
                                                        day: 'numeric',
                                                        hour: '2-digit',
                                                        minute: '2-digit'
                                                    }) : 'N/A',
                                                    bankLevel: sectionDetail.markers.bank_level,
                                                    level: m.level_cm,
                                                    danger: m.danger
                                                }))
                                            }
                                            margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                                        >
                                            <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
                                            <XAxis
                                                dataKey="time"
                                                tick={{ fontSize: 11 }}
                                            />
                                            <YAxis
                                                tick={{ fontSize: 12 }}
                                                domain={[0, sectionDetail.markers.bank_level + 1]}
                                                ticks={generateYAxisTicks(sectionDetail.markers.bank_level)}
                                                label={{ value: 'MSL', angle: -90, position: 'insideLeft' }}
                                            />
                                            <ReferenceLine
                                                y={bankLevel}
                                                stroke="#ef4444"
                                                strokeDasharray="4 4"
                                                strokeWidth={2}
                                                label={{
                                                    value: `Bank Level (${bankLevel} MSL)`,
                                                    position: 'insideTopRight',
                                                    fill: '#ef4444',
                                                    fontSize: 11
                                                }}
                                            />
                                            <Tooltip
                                                content={({ active, payload }) => {
                                                    if (active && payload && payload.length) {
                                                        const data = payload[0].payload;
                                                        return (
                                                            <div className="bg-white p-3 border border-gray-200 rounded-lg shadow-lg">
                                                                <p className="text-gray-600 text-sm mb-1">{data.fullDate}</p>
                                                                <p className="font-bold text-blue-600">
                                                                    ระดับน้ำ: {data.level} MSL
                                                                </p>
                                                            </div>
                                                        );
                                                    }
                                                    return null;
                                                }}
                                            />
                                            <Legend />
                                            <Line
                                                type="monotone"
                                                dataKey="level"
                                                stroke="#3b82f6"
                                                strokeWidth={2}
                                                dot={(props) => {
                                                    const { cx, cy, payload } = props as { cx?: number; cy?: number; payload?: { level: number } };
                                                    if (cx === undefined || cy === undefined || !payload) return null;
                                                    let color = '#3b82f6'; // blue - safe
                                                    if (payload.level > bankLevel) {
                                                        color = '#ef4444'; // red - danger
                                                    } else if (payload.level > bankLevel - 0.25) {
                                                        color = '#eab308'; // yellow - warning
                                                    }
                                                    return (
                                                        <circle
                                                            key={`dot-${cx}-${cy}`}
                                                            cx={cx}
                                                            cy={cy}
                                                            r={5}
                                                            fill={color}
                                                            stroke={color}
                                                            strokeWidth={2}
                                                        />
                                                    );
                                                }}
                                                activeDot={{ r: 7, fill: '#2563eb' }}
                                                name="Water Level (MSL)"
                                            />
                                            <Brush
                                                dataKey="time"
                                                height={30}
                                                stroke="#3b82f6"
                                                fill="#f0f9ff"
                                            />
                                        </LineChart>
                                    </ResponsiveContainer>
                                </div>
                            </div>

                            {/* {groupByDate(sectionDetail.markers.detail).map(([dateKey, markers]) => (
                                <div key={dateKey} className="bg-white rounded-xl shadow-lg overflow-hidden border border-gray-100">
                                    <div
                                        className="bg-gradient-to-r from-blue-600 to-cyan-600 text-white px-6 py-3 flex items-center justify-between cursor-pointer"
                                        onClick={() => setSortOrder(sortOrder === 'desc' ? 'asc' : 'desc')}
                                    >
                                        <h3 className="text-lg font-semibold flex items-center gap-2">
                                            {formatDateHeader(dateKey)}
                                        </h3>
                                        <span className="text-sm opacity-75 flex items-center gap-1">
                                            {markers.length} รายการ
                                            <span className="text-xs">{sortOrder === 'desc' ? '▼' : '▲'}</span>
                                        </span>
                                    </div>

                                    <div className="bg-gray-50 px-6 py-3 grid grid-cols-10 gap-4 text-sm font-semibold text-gray-600 border-b">
                                        <div className="col-span-2">เวลา</div>
                                        <div className="col-span-3">Water Level</div>
                                        <div className="col-span-2">Status</div>
                                        <div className="col-span-3">Flood</div>
                                    </div>

                                    <div className="divide-y divide-gray-100">
                                        {markers.map((marker: LocationDetail, index: number) => (
                                            <div
                                                key={`${marker.measured_at}-${index}`}
                                                className="px-6 py-4 grid grid-cols-10 gap-4 items-center hover:bg-blue-50 transition-colors"
                                            >
                                                <div className="col-span-2 text-gray-700 font-medium">
                                                    {formatTime(marker.measured_at)}
                                                </div>

                                                <div className="col-span-3">
                                                    <span className="text-2xl font-bold text-blue-600">
                                                        {marker.level_cm ? marker.level_cm.toFixed(2) : 'N/A'}
                                                    </span>
                                                    <span className="text-sm text-gray-500 ml-1">MSL</span>
                                                </div>

                                                <div className="col-span-2">
                                                    {marker.danger ? (
                                                        <span className={`px-3 py-1 rounded-full text-xs font-semibold border ${getDangerColor(marker.danger)}`}>
                                                            {marker.danger}
                                                        </span>
                                                    ) : (
                                                        <span className="text-gray-400">-</span>
                                                    )}
                                                </div>

                                                <div className="col-span-3">
                                                    {marker.is_flooded !== null ? (
                                                        <span className={`px-3 py-1 rounded-full text-xs font-semibold ${marker.is_flooded ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'}`}>
                                                            {marker.is_flooded ? 'Flooded' : 'Normal'}
                                                        </span>
                                                    ) : (
                                                        <span className="text-gray-400">-</span>
                                                    )}
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            ))} */}
                        </div>
                    )}
                </div>
            ) : (
                <div className="flex items-center justify-center min-h-[60vh]">
                    <div className="text-center">
                        <div className="animate-spin rounded-full h-16 w-16 border-b-4 border-blue-600 mx-auto mb-4"></div>
                        <p className="text-gray-600 text-lg">Loading water level data...</p>
                    </div>
                </div>
            )}

            <Footer />
        </div>
    );
};

export default SectionDetail;