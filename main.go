package main

import(
    "context"
    "log"
    "net"
    "sync"
    //import the generated protobuf code
    pb "github.com/jojojolin/learn_microservices/proto/consignment"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

const (
    port = ":50051"
)

type repository interface {
    Create(*pb.Consignment) (*pb.Consignment, error)
    GetAll() []*pb.Consignment
}

//Repository -fake
type Repository struct {
    mu sync.RWMutex
    consignments []*pb.Consignment
}

func (repo *Repository)GetAll() []*pb.Consignment{
    return repo.consignments
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest)(*pb.Response, error){
    consignments := s.repo.GetAll()
    return &pb.Response{Consignments: consignments},nil

}

//Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
    repo.mu.Lock()
    updated := append(repo.consignments, consignment)
    repo.consignments = updated
    repo.mu.Unlock()
    return consignment, nil
}

type service struct{
    repo repository
}

func(s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error){
    // Save our consignment
    consignment, err :=s.repo.Create(req)
    if err!=nil {
        return nil, err
    }
    //Return matching the `Response` message we created in our protobuf def
    return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
    repo := &Repository{}
    
    //Set-up our gRPC server
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s:=grpc.NewServer()
    
    //Register our service with the gRPC server, this will tie our implementation into the auto-generated interface code for our protobuf definition.
    pb.RegisterShippingServiceServer(s, &service{repo})
    
    //Register reflection service on gRPC server.
    reflection.Register(s)
    
    log.Println("Running on port:",port)
    if err :=s.Serve(lis); err !=nil { //%v is the value in default format
        log.Fatalf("failed to serve: %v",err)
    }
}

