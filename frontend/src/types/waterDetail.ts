export interface LocationDetail {
    level_cm: number | null;
    image: string;
    danger: string | null;
    is_flooded: boolean | null;
    source: {
        String: string;
        Valid: boolean;
    } | null;
    measured_at: string | null;
    note: string | null;
}

export interface WaterLocationDetail {
    station_id: number;
    bank_level: number;
    detail: LocationDetail[];
}

export interface WaterDetailResponse {
    markers: WaterLocationDetail;
}




