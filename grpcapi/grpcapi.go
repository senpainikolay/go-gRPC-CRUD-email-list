package grpcapi

import (
	"context"
	"database/sql"
	"log"
	"net"
	"simple_proj/db_sql"
	pb "simple_proj/proto"
	"time"

	"google.golang.org/grpc"
)

type MailServer struct {
	pb.UnimplementedMailingListServiceServer
	db *sql.DB
}

func pbEntryToDbEntry(pbEntry *pb.EmailEntry) db_sql.EmailEntry {
	t := time.Unix(pbEntry.ConfirmedAt, 0)
	return db_sql.EmailEntry{
		Id:          pbEntry.Id,
		Email:       pbEntry.Email,
		ConfirmedAt: &t,
		OptOut:      pbEntry.OptOut,
	}
}

func dbEntryToPbEntry(mdbEntry *db_sql.EmailEntry) pb.EmailEntry {
	return pb.EmailEntry{
		Id:          mdbEntry.Id,
		Email:       mdbEntry.Email,
		ConfirmedAt: mdbEntry.ConfirmedAt.Unix(),
		OptOut:      mdbEntry.OptOut,
	}
}

func emailResponse(db *sql.DB, email string) (*pb.EmailResponse, error) {
	entry, err := db_sql.GetEmail(db, email)
	if err != nil {
		return &pb.EmailResponse{}, err
	}
	if entry == nil {
		return &pb.EmailResponse{}, nil
	}

	res := dbEntryToPbEntry(entry)

	return &pb.EmailResponse{EmailEntry: &res}, nil
}

func (s *MailServer) GetEmail(ctx context.Context, req *pb.GetEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC GetEmail: %v\n", req)
	return emailResponse(s.db, req.EmailAddr)
}

func (s *MailServer) GetEmailBatch(ctx context.Context, req *pb.GetEmailBatchRequest) (*pb.GetEmailBatchResponse, error) {
	log.Printf("gRPC GetEmailBatch: %v\n", req)
	params := db_sql.GetEmailBatchQueryParams{
		Page:  int(req.Page),
		Count: int(req.Count),
	}

	dbEntries, err := db_sql.GetEmailBatch(s.db, params)
	if err != nil {
		return &pb.GetEmailBatchResponse{}, err
	}

	pbEntries := make([]*pb.EmailEntry, 0, len(dbEntries))
	for i := 0; i < len(dbEntries); i++ {
		entry := dbEntryToPbEntry(&dbEntries[i])
		pbEntries = append(pbEntries, &entry)
	}

	return &pb.GetEmailBatchResponse{EmailEntries: pbEntries}, nil
}

func (s *MailServer) CreateEmail(ctx context.Context, req *pb.CreateEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC CreateEmail: %v\n", req)

	err := db_sql.CreateEmail(s.db, req.EmailAddr)
	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, req.EmailAddr)
}

func (s *MailServer) UpdateEmail(ctx context.Context, req *pb.UpdateEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC UpdateEmail: %v\n", req)

	entry := pbEntryToDbEntry(req.EmailEntry)

	err := db_sql.UpdateEmail(s.db, entry)
	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, entry.Email)
}

func (s *MailServer) DeleteEmail(ctx context.Context, req *pb.DeleteEmailRequest) (*pb.EmailResponse, error) {
	log.Printf("gRPC DeleteEmail: %v\n", req)

	err := db_sql.DeleteEmail(s.db, req.EmailAddr)
	if err != nil {
		return &pb.EmailResponse{}, err
	}

	return emailResponse(s.db, req.EmailAddr)
}

func Serve(db *sql.DB, bind string) {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		log.Fatalf("gRPC server error: failure to bind %v\n", bind)
	}

	grpcServer := grpc.NewServer()

	mailServer := MailServer{db: db}

	pb.RegisterMailingListServiceServer(grpcServer, &mailServer)

	log.Printf("gRPC API server listening on %v\n", bind)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("gRPC server error: %v\n", err)
	}
}
