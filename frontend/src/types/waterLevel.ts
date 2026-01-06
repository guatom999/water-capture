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

interface MarkersWithLocation {
    LocationID: string,
    LocationName: string,
    LocationDescription: string,
    Latitude: number,
    Longitude: number,
    IsActive: boolean,
    WaterLevelID: string,
    LevelCm: number,
    Image: string,
    Danger: string,
    IsFlooded: boolean,
    MeasuredAt: string,
    Note: string
}