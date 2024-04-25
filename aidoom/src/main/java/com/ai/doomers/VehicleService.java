package com.ai.doomers;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.util.HashMap;
import java.util.List;
import java.util.Scanner;
import java.util.stream.Collectors;

public class VehicleService {

    private static final String CSV_FILE = "Car_Inventory.csv";
    private HashMap<Integer, Vehicle> activeVehicles = new HashMap<>(); // Vehicles that match the criteria
    private HashMap<Integer, Vehicle> inactiveVehicles = new HashMap<>(); //Vehicles that are not in use
    private HashMap<Integer, Vehicle> inUseVehicles = new HashMap<>(); //Vehicles that are in use

    //Loading Vehicle Data from CSV.
    public void inventoryProcessor() {
        try (FileInputStream vehicleData = new FileInputStream(CSV_FILE); Scanner inFs = new Scanner(vehicleData)) {
            while (inFs.hasNextLine()) {
                String line = inFs.nextLine();
                String[] vData = line.split(",");
                
                // 0-model, 1-make, 2-capacity, 3-wheelchairADA, 4-visionImpairedADA, 5-id
                int id = Integer.parseInt(vData[5]);
                Vehicle vehicle = new Vehicle(vData[0], vData[1], Integer.parseInt(vData[2]), Boolean.parseBoolean(vData[3]), Boolean.parseBoolean(vData[4]), id);
                //Adds vehicles to the activeVehicles Map
                inactiveVehicles.put(id, vehicle); 
            }
        } catch (FileNotFoundException e) {
            System.out.println("The file was not found: " + e.getMessage());
        } catch (IOException e) {
            System.out.println("An error occurred while closing the file: " + e.getMessage());
        }
    }

    //Return a list of vehicles that match the criteria 
    public List<Vehicle> findVehicles(int capacity, boolean wheelchairADA, boolean visionImpairedADA) {
       return inactiveVehicles.values().stream()
               .filter(v -> v.getCapacity() >= capacity && v.isWheelchairADA() == wheelchairADA && v.isVisionImpairedADA() == visionImpairedADA)
               .collect(Collectors.toList());
    }

    //Removes the Vehicles from the Act
    public void markVehicleAsUsed(int id) {
        Vehicle vehicle = activeVehicles.remove(id);
        if (vehicle != null) {
            inUseVehicles.put(id, vehicle);
        }
    }

    public void returningVehicles(int id) {
        Vehicle vehicle =  inUseVehicles.remove(id);
        if (vehicle != null) {
            inactiveVehicles.put(id, vehicle);
        } 
    }

 
    public void activateVehicle(int id) {
        Vehicle vehicle = inactiveVehicles.remove(id);
        if (vehicle != null) {
            activeVehicles.put(id, vehicle);
        }
    }
}
