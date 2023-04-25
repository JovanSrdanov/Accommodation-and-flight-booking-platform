import yaml
from yaml.loader import SafeLoader
import sys, getopt

#Because...
true = True
false = False

consoleOut = sys.stdout

def  getParameters():

    #Defaults
    inputFile = "docker-compose.yml"
    outputFile = "docker-compose.local.yml"
    envFile = ".env.local"

    serviceName = ""

    # options - stores all options with its parameters (option, parameter)
    # first parameter is argument list from terminal
    # second parameter is allowed options, if options espect argument they are followed by ":"
    # third parameter can be added for long options (--option)
    options, arguments = getopt.getopt(sys.argv[1:],"i:ho:e:s:")


    for option, argument in options:
        if option == "-h":
            print("LocalDevStartup.py [-i <inputFile -o <outputFile> -e <envFile>] <service-name>")
            sys.exit()
        elif option  == "-i":
            inputFile = argument
        elif option == "-o":
            outputFile = argument
        elif option == "-e":
            envFile = argument

    if len(arguments) == 0:
        print("Service name not provided, enter -h for help")
        sys.exit()

    serviceName = arguments[0]

    exit = false
    if envFile == "":
        print("Environment file not provided")
        exit = true
    if inputFile == "":
        print("Input file not provided")
        exit = true
    if outputFile == "":
        print("Output file not provided")
        exit = true

    if(exit):
        sys.exit()
    
    return inputFile, outputFile, envFile, serviceName




# Generates .env file for local development
# Changes all names of services  in links to localhost
# Returns which service should expose which port so our service can connect to it when run locally
def generateEnvFile(service, envFile, serviceNames):
    envVars = service["environment"]

    with open(envFile, "w") as sys.stdout:
        exposingServicePorts = {}
        for var in envVars:


            for serviceName in serviceNames:
                #For example mongo:27017
                portPrefix = serviceName + ":"

                found = var.find(portPrefix)
                if found != -1:
                    port = var[found + len(portPrefix):]
                    # so, for mongo service port 27017 should be opened for local development purposes
                    exposingServicePorts[serviceName] = port
                    # in local we are targeting localhost not container name
                    var = var.replace(portPrefix, "localhost:")


            line = "export " + var
            print(line)
    # Sets stdout again to console output
    sys.stdout = consoleOut
    print("Environment variables written in '" + envFile + "' file")
    return exposingServicePorts



def exposePorts(exposingServicePorts, services):
    for service in exposingServicePorts:
        port = exposingServicePorts[service]
        # Map local port to same container port (this is the only con if that port is already taken
        # in that case you will have to assign it manually)
        mapping = port + ":" + port

        #if there is no ports mapping already defined for service make ports label
        if "ports" not in services[service]:
            services[service]["ports"] = []

        services[service]["ports"].append(mapping)



def main():
    inputFile, outputFile, envFile, serviceName = getParameters()

    with open(inputFile, "r") as file:
        data = yaml.load(file, Loader = SafeLoader)
        services = data["services"]

        if serviceName not in services.keys():
            print("There is no service '" + serviceName + "' in passed docker compose file")
            sys.exit()

        targetService = services[serviceName]

        print("Generating files for local development...")

        #exposingServicePorts map with content: {service : port}
        exposingServicePorts = generateEnvFile(targetService, envFile, list(services.keys()))
        exposePorts(exposingServicePorts, services)

        # Removing target service
        services.pop(serviceName)
        data["services"] = services
        
        
        with open(outputFile, 'w') as f:
            data = yaml.dump(data, f, sort_keys=False, default_flow_style=False)
        print("Docker compose written in '" + outputFile + "' file")



if __name__ == "__main__":
    main()