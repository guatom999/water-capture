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

                        {mapData.markers.map((marker: MarkersWithLocation) => (
                            <Marker key={marker.location_id} position={[marker.latitude, marker.longitude]}>
                                <Popup>
                                    <Link to={`/markers/detail?location_id=${marker.location_id}`}>
                                        <div className="font-bold">{marker.location_name}</div>
                                    </Link>
                                    <div className="whitespace-pre-line">{marker.note}</div>
                                    <div className="text-sm text-gray-600 mt-2">
                                        Level: {marker.level_cm} cm
                                    </div>
                                    <img src={marker.image} alt="" />
                                </Popup>
                            </Marker>
                        ))}
                    </MapContainer>
                </div>
            </div>
        </section>
    );
};

export default MapSection;