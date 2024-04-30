package com.ai.doomers;

import java.util.HashSet;
import java.util.Set;

import com.fasterxml.jackson.annotation.JsonValue;

public enum VehicleFlag {
    WheelchairAccess,
    VisionAccess;

    public static VehicleFlag parse(String raw) {
        switch(raw) {
        case "wheelchair": return WheelchairAccess;
        case "vision": return VisionAccess;
        }

        throw new IllegalArgumentException(raw + " is not a valid flag");
    }

    public static Set<VehicleFlag> parseCommaSeperated(String raw) {
        Set<VehicleFlag> out = new HashSet<>();
        if(raw == null) {
            return out;
        }

        for(String s: raw.split(",")) {
            out.add(parse(s));
        }
        return out;
    }

    public String toString() {
        switch(this) {
        case WheelchairAccess: return "wheelchair";
        case VisionAccess: return "vision";
        }

        throw new Error("toString was not exhaustive");
    }

    @JsonValue
    public String getValue() {
        return toString();
    }
}
