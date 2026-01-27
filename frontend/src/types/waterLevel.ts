// export interface MapMarkerResponse {
//     center: Center;
//     markers: Markers[];
//     zoom: number;
// }

export interface MapMarkerResponse {
    markers: MarkersWithLocation[];
}

// interface Center {
//     latitude: number;
//     longitude: number;
//     name: string;
// }

// interface Markers {
//     id: string;
//     name: string;
//     description: string;
//     latitude: number;
//     longitude: number;
//     level: number;
//     icon: string;
//     timestamp: string;
// }

export interface MarkersWithLocation {
    location_id: string,
    location_name: string,
    location_description: string,
    latitude: number,
    longitude: number,
    is_active: boolean,
    bank_level: number,
    water_level_id: string,
    level_cm: number,
    image: string,
    danger: string,
    is_flooded: boolean,
    measured_at: string,
    note: string
}