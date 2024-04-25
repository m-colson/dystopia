package com.ai.doomers;


public class Vehicle {

    private String model;
    private String make;
    private int capacity;
    private int id;
    private boolean wheelchairADA;
    private boolean visionImpairedADA;


    // Modified constructor to generate a random ID if not provided
    public Vehicle(String carModel, String carMake, int carCapacity, boolean wheelchair, boolean visionImpaired, int carId) {
        this.model = carModel;
        this.make = carMake;
        this.capacity = carCapacity;
        this.visionImpairedADA = visionImpaired ;
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

    public boolean isVisionImpairedADA() {
        return visionImpairedADA;
    }
    public boolean isWheelchairADA() {
        return wheelchairADA;
    }
}
