package com.ai.doomers;

import java.util.HashSet;
import java.util.Set;

public class Vehicle {

    private String model;
    private String make;
    private int capacity;
    private int id;
    private Set<VehicleFlag> flags;


    // Modified constructor to generate a random ID if not provided
    public Vehicle(String carModel, String carMake, int carCapacity, boolean wheelchair, boolean visionImpaired, int carId) {
        this.model = carModel;
        this.make = carMake;
        this.capacity = carCapacity;
        
        this.flags = new HashSet<>();
        if(wheelchair) {
            flags.add(VehicleFlag.WheelchairAccess);
        }
        if(visionImpaired) {
            flags.add(VehicleFlag.VisionAccess);
        }

        this.id = carId;
    }

    // Getters and setters as before
    public String getModel() {
        return model;
    }

    public String getMake() {
        return make;
    }

    public int getCapacity() {
        return capacity;
    }

    public int getId() {
        return id;
    }

    public Set<VehicleFlag> getFlags() {
        return flags;
    }
}
