import { useState, useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import { Link } from 'react-router-dom';
import 'leaflet/dist/leaflet.css';
import L, { Handler } from 'leaflet';
import icon from 'leaflet/dist/images/marker-icon.png';
import iconShadow from 'leaflet/dist/images/marker-shadow.png';
import Header from './Header';
import { getMapMarkersService } from '../services/waterLevelService';
import type { MapMarkerResponse } from '../types/waterLevel';

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
        <section className="bg-gray-50 w-full h-screen relative">
            <div className='w-full h-full'>
                <MapContainer
                    // center={[mapData.Markers[0].Latitude, mapData.Markers[0].Longitude]}
                    center={[13.7567, 100.5115]}
                    zoom={10}
                    scrollWheelZoom={true}
                    style={{ height: '100%', width: '100%' }}
                >
                    <TileLayer
                        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    />

                    {mapData.markers.map((marker: any) => (
                        <Marker key={marker.LocationID} position={[marker.Latitude, marker.Longitude]}>
                            <Popup>
                                <Link to={`/markers/detail?location_id=${marker.LocationID}`}>
                                    <div className="font-bold">{marker.LocationName}</div>
                                </Link>
                                <div className="whitespace-pre-line">{marker.Note}</div>
                                <div className="text-sm text-gray-600 mt-2">
                                    Level: {marker.LevelCm} cm
                                </div>
                                <img src={marker.Image} alt="" />
                            </Popup>
                        </Marker>
                    ))}
                </MapContainer>
            </div>
        </section>
    );
};

export default MapSection;