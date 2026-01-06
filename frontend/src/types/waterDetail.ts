export interface LocationDetail {
    location_id: number;
    level_cm: number | null;
    image: string ;
    danger: string | null;
    is_flooded: boolean | null;
    source: {
        String: string;
        Valid: boolean;
    } | null;
    measured_at: string | null;
    note: string | null;
}

export interface WaterDetailResponse {
    markers: LocationDetail[];
}




