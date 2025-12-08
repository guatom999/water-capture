export interface MapMarkerResponse {
    center: Center;
    markers: Markers[];
    zoom: number;
}

interface Center {
    latitude: number;
    longitude: number;
    name: string;
}

interface Markers {
    id: string;
    name: string;
    description: string;
    latitude: number;
    longitude: number;
    level: number;
    icon: string;
    timestamp: string;
}