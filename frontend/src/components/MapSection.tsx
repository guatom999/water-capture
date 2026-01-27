import { useState, useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import { Link } from 'react-router-dom';
import 'leaflet/dist/leaflet.css';
import { getMapMarkersService } from '../services/waterLevelService';
import type { MapMarkerResponse, MarkersWithLocation } from '../types/waterLevel';

// let DefaultIcon = L.icon({
//     iconUrl: icon,
//     shadowUrl: iconShadow,
//     iconSize: [25, 41],
//     iconAnchor: [12, 41]
// });

// L.Marker.prototype.options.icon = DefaultIcon;

const MapSection = () => {
    const [mapData, setMapData] = useState<MapMarkerResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchMapData = async () => {
            try {
                setLoading(true);
                const data = await getMapMarkersService();
                setMapData(data as MapMarkerResponse);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to fetch map data');
                console.error('Error fetching map data:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchMapData();
    }, []);

    if (loading) {
        return (
            <section className="bg-gray-50 w-full h-screen flex items-center justify-center">
                <div className="text-xl">Loading map data...</div>
            </section>
        );
    }

    if (error || !mapData) {
        return (
            <section className="bg-gray-50 w-full h-screen flex items-center justify-center">
                <div className="text-xl text-red-500">Error: {error || 'No data available'}</div>
            </section>
        );
    }

    return (
        <section className="bg-gradient-to-b from-gray-50 to-gray-100 w-full py-8 px-4 sm:px-6 lg:px-8">
            <div className="max-w-7xl mx-auto">
                <div className="mb-6 text-center">
                    <h2 className="text-2xl sm:text-3xl font-bold text-gray-800 mb-2">
                        แผนที่ระดับน้ำ
                    </h2>
                    <p className="text-gray-600">
                        คลิกที่ marker เพื่อดูรายละเอียดระดับน้ำในแต่ละจุด
                    </p>
                </div>

                {/* Map Container */}
                <div className="rounded-2xl overflow-hidden shadow-xl border border-gray-200" style={{ height: '70vh' }}>
                    <MapContainer
                        center={[13.7567, 100.5115]}
                        zoom={8}
                        scrollWheelZoom={true}
                        style={{ height: '100%', width: '100%' }}
                    >
                        <TileLayer
                            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        />

                        {mapData.markers.map((marker: MarkersWithLocation) => {
                            // Helper functions for this marker
                            const getDangerInfo = (danger: string | null) => {
                                switch (danger?.toUpperCase()) {
                                    case 'CRITICAL':
                                        return { bgColor: 'bg-red-500', text: 'วิกฤต' };
                                    case 'DANGER':
                                        return { bgColor: 'bg-orange-500', text: 'อันตราย' };
                                    case 'WATCH':
                                        return { bgColor: 'bg-yellow-500', text: 'เฝ้าระวัง' };
                                    case 'SAFE':
                                        return { bgColor: 'bg-green-500', text: 'ปลอดภัย' };
                                    default:
                                        return { bgColor: 'bg-gray-400', text: 'ไม่ทราบ' };
                                }
                            };

                            const formatTime = (dateString: string | null) => {
                                if (!dateString) return 'N/A';
                                return new Date(dateString).toLocaleString('th-TH', {
                                    hour: '2-digit',
                                    minute: '2-digit'
                                });
                            };

                            const dangerInfo = getDangerInfo(marker.danger);

                            return (
                                <Marker key={marker.location_id} position={[marker.latitude, marker.longitude]}>
                                    <Popup className="custom-popup" minWidth={260} maxWidth={300}>
                                        <div className="p-0 -m-3">
                                            {/* Header */}
                                            <div className="p-3 border-b border-gray-200">
                                                <div className="font-bold text-gray-800">{marker.location_name}</div>
                                                <div className="flex items-center gap-1 mt-1">
                                                    <span className={`w-2 h-2 rounded-full ${dangerInfo.bgColor}`}></span>
                                                    <span className="text-xs text-gray-600">{dangerInfo.text}</span>
                                                </div>
                                            </div>

                                            {/* Water Level Info */}
                                            <div className="p-3">
                                                <div className="text-center mb-3">
                                                    <div className="text-2xl font-bold text-blue-600">
                                                        {(marker.level_cm / 100).toFixed(2)} m
                                                    </div>
                                                    <div className="text-xs text-gray-500">
                                                        ระดับตลิ่ง: {marker.bank_level?.toFixed(2) || '?'} m
                                                    </div>
                                                </div>

                                                {/* Status row */}
                                                <div className="flex items-center justify-between text-xs text-gray-500 mb-3">
                                                    <span>อัพเดท: {formatTime(marker.measured_at)}</span>
                                                    <span className={`px-2 py-0.5 rounded text-xs font-medium ${marker.is_flooded ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                                                        {marker.is_flooded ? 'น้ำท่วม' : 'ปกติ'}
                                                    </span>
                                                </div>

                                                {/* Action Button */}
                                                <Link
                                                    to={`/markers/detail?location_id=${marker.location_id}`}
                                                    className="block w-full text-center bg-blue-500 text-white py-2 px-4 rounded text-sm font-medium hover:bg-blue-600 transition-colors"
                                                    style={{
                                                        color: 'white'
                                                    }}
                                                >
                                                    ดูรายละเอียด
                                                </Link>
                                            </div>
                                        </div>
                                    </Popup>
                                </Marker>
                            );
                        })}
                    </MapContainer>
                </div>
            </div>
        </section>
    );
};

export default MapSection;