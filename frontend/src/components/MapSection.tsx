import { useState, useEffect, useCallback, useRef } from 'react';
import { MapContainer, TileLayer, Marker, Popup, GeoJSON, useMap, CircleMarker } from 'react-leaflet';
import { Link } from 'react-router-dom';
import type { Layer, LeafletMouseEvent } from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { getMapMarkersService } from '../services/waterLevelService';
import type { MapMarkerResponse, MarkersWithLocation } from '../types/waterLevel';
import { pathumthaniGeoJSON } from '../data/pathumthani';
import { bangkokGeoJSON } from '../data/bangkok';

// Styles for province boundaries
const borderStyle = {
    color: '#3388ff',
    weight: 2,
    fillColor: '#3388ff',
    fillOpacity: 0.1,
};

const highlightStyle = {
    color: '#ff6b35',
    weight: 4,
    fillColor: '#ff6b35',
    fillOpacity: 0.25,
};

// Get color based on danger level
const getDangerColor = (danger: string | null): string => {
    switch (danger?.toUpperCase()) {
        case 'CRITICAL':
            return '#ef4444'; // red
        case 'DANGER':
            return '#f97316'; // orange
        case 'WATCH':
            return '#eab308'; // yellow
        case 'SAFE':
            return '#22c55e'; // green
        default:
            return '#9ca3af'; // gray
    }
};

// Province Layer Component with zoom on click and station circles
interface ProvinceLayerProps {
    data: GeoJSON.FeatureCollection;
    provinceId: number;
    markers: MarkersWithLocation[];
    selectedProvinceId: number | null;
    onProvinceClick: (provinceId: number) => void;
    zoomLevel: number;
}

const ProvinceLayer = ({ data, provinceId, markers, selectedProvinceId, onProvinceClick, zoomLevel }: ProvinceLayerProps) => {
    const map = useMap();

    // Filter markers for this province
    const provinceMarkers = markers.filter(marker => marker.province_id === provinceId);

    const onEachFeature = (feature: GeoJSON.Feature, layer: Layer) => {
        const props = feature.properties;
        // Only show tooltip when zoom level is 10 or less
        if (props && zoomLevel <= 10) {
            layer.bindTooltip(
                `<div style="padding: 8px; min-width: 150px;">
                    <div style="font-weight: bold; font-size: 14px; margin-bottom: 4px; color: #1f2937;">
                        🏛️ ${props.NAME_1}
                    </div>
                    <div style="font-size: 12px; color: #6b7280;">
                        ประเทศ: ${props.NAME_0}
                    </div>
                    <div style="font-size: 12px; color: #6b7280;">
                        รหัสจังหวัด: ${props.ID_1}
                    </div>
                    <div style="font-size: 12px; color: #3b82f6; margin-top: 4px;">
                        📍 สถานี: ${provinceMarkers.length} แห่ง
                    </div>
                </div>`,
                {
                    sticky: true,
                    direction: 'top',
                    className: 'province-tooltip',
                    offset: [0, -10]
                }
            );
        }

        // Add hover and click events
        layer.on({
            mouseover: (e: LeafletMouseEvent) => {
                const target = e.target;
                // Don't highlight if another province is selected (this one is faded out)
                if (selectedProvinceId !== null && selectedProvinceId !== provinceId) {
                    return;
                }
                target.setStyle(highlightStyle);
                target.bringToFront();
            },
            mouseout: (e: LeafletMouseEvent) => {
                const target = e.target;
                // Don't reset if another province is selected (keep faded out style)
                if (selectedProvinceId !== null && selectedProvinceId !== provinceId) {
                    return;
                }
                target.setStyle(borderStyle);
            },
            click: (e: LeafletMouseEvent) => {
                const bounds = e.target.getBounds();
                map.fitBounds(bounds, { padding: [50, 50], maxZoom: 11 });
                onProvinceClick(provinceId);
            }
        });
    };

    // Adjust style based on selection state
    const getStyle = () => {
        if (selectedProvinceId === provinceId) {
            // Selected province - keep fill but disable interaction for marker clicks
            return { ...borderStyle, interactive: false };
        } else if (selectedProvinceId !== null) {
            // Other province when one is selected - fade out
            return { ...borderStyle, fillOpacity: 0 };
        }
        // No selection - normal style
        return borderStyle;
    };
    const currentStyle = getStyle();

    return (
        <>
            <GeoJSON
                key={`province-${provinceId}-${selectedProvinceId === null ? 'none' : selectedProvinceId}`}
                data={data}
                style={currentStyle}
                onEachFeature={onEachFeature}
            />
            {/* Show station circles when this province is selected */}
            {selectedProvinceId === provinceId && provinceMarkers.map((marker) => (
                <CircleMarker
                    key={marker.station_id}
                    center={[marker.latitude, marker.longitude]}
                    radius={12}
                    pathOptions={{
                        color: getDangerColor(marker.danger),
                        fillColor: getDangerColor(marker.danger),
                        fillOpacity: 0.7,
                        weight: 2
                    }}
                >
                    <Popup className="custom-popup" minWidth={260} maxWidth={300}>
                        <div className="p-0 -m-3">
                            <div className="p-3 border-b border-gray-200">
                                <div className="font-bold text-gray-800">{marker.location_name}</div>
                                <div className="flex items-center gap-1 mt-1">
                                    <span
                                        className="w-2 h-2 rounded-full"
                                        style={{ backgroundColor: getDangerColor(marker.danger) }}
                                    ></span>
                                    <span className="text-xs text-gray-600">
                                        {marker.danger || 'ไม่ทราบ'}
                                    </span>
                                </div>
                            </div>

                            <div className="p-3">
                                <div className="text-center mb-3">
                                    <div className="text-2xl font-bold text-blue-600">
                                        {marker.level_cm} msl
                                    </div>
                                    <div className="text-xs text-gray-500">
                                        ระดับตลิ่ง: {marker.bank_level?.toFixed(2) || '?'} msl
                                    </div>
                                </div>

                                <div className="flex items-center justify-between text-xs text-gray-500 mb-3">
                                    <span>อัพเดท: {marker.measured_at ? new Date(marker.measured_at).toLocaleString('th-TH', {
                                        hour: '2-digit',
                                        minute: '2-digit'
                                    }) : 'N/A'}</span>
                                    <span className={`px-2 py-0.5 rounded text-xs font-medium ${marker.is_flooded ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                                        {marker.is_flooded ? 'น้ำท่วม' : 'ปกติ'}
                                    </span>
                                </div>

                                <Link
                                    to={`/markers/detail?station_id=${marker.station_id}`}
                                    className="block w-full text-center bg-blue-500 text-white py-2 px-4 rounded text-sm font-medium hover:bg-blue-600 transition-colors"
                                    style={{ color: 'white' }}
                                >
                                    ดูรายละเอียด
                                </Link>
                            </div>
                        </div>
                    </Popup>
                </CircleMarker>
            ))}
        </>
    );
};

// Component to handle zoom changes
interface ZoomHandlerProps {
    onZoomChange: (zoom: number) => void;
}

const ZoomHandler = ({ onZoomChange }: ZoomHandlerProps) => {
    const map = useMap();
    const initializedRef = useRef(false);

    useEffect(() => {
        const handleZoom = () => {
            onZoomChange(map.getZoom());
        };

        map.on('zoomend', handleZoom);
        // Only set initial zoom once on mount
        if (!initializedRef.current) {
            initializedRef.current = true;
            onZoomChange(map.getZoom());
        }

        return () => {
            map.off('zoomend', handleZoom);
        };
    }, [map, onZoomChange]);

    return null;
};

const MapSection = () => {
    const [zoomLevel, setZoomLevel] = useState<number>(8);
    const [mapData, setMapData] = useState<MapMarkerResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedProvinceId, setSelectedProvinceId] = useState<number | null>(null);

    const handleProvinceClick = (provinceId: number) => {
        // Toggle: if clicking same province, deselect. Otherwise select new province
        setSelectedProvinceId(prev => prev === provinceId ? null : provinceId);
    };

    // Handle zoom changes - reset province selection when zoom is below 10
    const handleZoomChange = useCallback((zoom: number) => {
        setZoomLevel(zoom);
        if (zoom < 10) {
            setSelectedProvinceId(null);
        }
    }, []);

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
                        {selectedProvinceId}
                    </h2>
                    <p className="text-gray-600">
                        คลิกที่ marker เพื่อดูรายละเอียดระดับน้ำในแต่ละจุด
                    </p>
                </div>

                {/* Map Container */}
                <div className="rounded-2xl overflow-hidden shadow-xl border border-gray-200" style={{ height: '70vh' }}>
                    <MapContainer
                        center={[13.7567, 100.5115]}
                        zoom={zoomLevel}
                        scrollWheelZoom={true}
                        style={{ height: '100%', width: '100%' }}
                        maxBounds={[[5.5, 97.0], [20.5, 106.0]]}  // ขอบเขตประเทศไทย
                        minZoom={5}
                        maxZoom={18}
                    >
                        <TileLayer
                            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                        />

                        <ZoomHandler onZoomChange={handleZoomChange} />

                        {/* Province Borders with zoom on click */}
                        <ProvinceLayer
                            data={pathumthaniGeoJSON as unknown as GeoJSON.FeatureCollection}
                            provinceId={13}
                            markers={mapData.markers}
                            selectedProvinceId={selectedProvinceId}
                            onProvinceClick={handleProvinceClick}
                            zoomLevel={zoomLevel}
                        />
                        <ProvinceLayer
                            data={bangkokGeoJSON as unknown as GeoJSON.FeatureCollection}
                            provinceId={30}
                            markers={mapData.markers}
                            selectedProvinceId={selectedProvinceId}
                            onProvinceClick={handleProvinceClick}
                            zoomLevel={zoomLevel}
                        />

                        {/* {mapData.markers.map((marker: MarkersWithLocation) => {
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
                                            <div className="p-3 border-b border-gray-200">
                                                <div className="font-bold text-gray-800">{marker.location_name}</div>
                                                <div className="flex items-center gap-1 mt-1">
                                                    <span className={`w-2 h-2 rounded-full ${dangerInfo.bgColor}`}></span>
                                                    <span className="text-xs text-gray-600">{dangerInfo.text}</span>
                                                </div>
                                            </div>

                                            <div className="p-3">
                                                <div className="text-center mb-3">
                                                    <div className="text-2xl font-bold text-blue-600">
                                                        {marker.level_cm} msl
                                                    </div>
                                                    <div className="text-xs text-gray-500">
                                                        ระดับตลิ่ง: {marker.bank_level?.toFixed(2) || '?'} msl
                                                    </div>
                                                </div>

                                                <div className="flex items-center justify-between text-xs text-gray-500 mb-3">
                                                    <span>อัพเดท: {formatTime(marker.measured_at)}</span>
                                                    <span className={`px-2 py-0.5 rounded text-xs font-medium ${marker.is_flooded ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'}`}>
                                                        {marker.is_flooded ? 'น้ำท่วม' : 'ปกติ'}
                                                    </span>
                                                </div>

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
                        })} */}

                    </MapContainer>
                </div>
            </div>
        </section>
    );
};

export default MapSection;