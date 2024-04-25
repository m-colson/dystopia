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
    private HashMap<Integer, Vehicle> activeVehicles = new HashMap<>();
    private HashMap<Integer, Vehicle> usedVehicles = new HashMap<>();

    public void inventoryProcessor() {
        try (FileInputStream vehicleData = new FileInputStream(CSV_FILE); Scanner inFs = new Scanner(vehicleData)) {
            while (inFs.hasNextLine()) {
                String line = inFs.nextLine();
                String[] vData = line.split(",");
                
                // 0-model, 1-make, 2-capacity, 3-wheelchairADA, 4-visionImpairedADA, 5-id
                int id = Integer.parseInt(vData[5]);
                Vehicle vehicle = new Vehicle(vData[0], vData[1], Integer.parseInt(vData[2]),
                                              Boolean.parseBoolean(vData[3]), Boolean.parseBoolean(vData[4]), id);
                activeVehicles.put(id, vehicle);
            }
        } catch (FileNotFoundException e) {
            System.out.println("The file was not found: " + e.getMessage());
        } catch (IOException e) {
            System.out.println("An error occurred while closing the file: " + e.getMessage());
        }
    }

    public List<Vehicle> findVehicles(int capacity, boolean wheelchairADA, boolean visionImpairedADA) {
        return activeVehicles.values().stream()
                .filter(v -> v.getCapacity() == capacity && v.isWheelchairADA() == wheelchairADA && v.isVisionImpairedADA() == visionImpairedADA)
                .collect(Collectors.toList());
    }

    public void markVehicleAsUsed(int id) {
        Vehicle vehicle = activeVehicles.remove(id);
        if (vehicle != null) {
            usedVehicles.put(id, vehicle);
        }
    }
}
