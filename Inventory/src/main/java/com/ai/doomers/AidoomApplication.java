package com.ai.doomers;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
@ComponentScan(basePackages = {"com.ai.doomers", "any.other.packages"})
public class AidoomApplication {

	public static void main(String[] args) {
		SpringApplication.run(AidoomApplication.class, args);
	}

}
